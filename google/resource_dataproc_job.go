package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"google.golang.org/api/dataproc/v1"
)

func resourceDataprocJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataprocJobCreate,
		Update: resourceDataprocJobUpdate,
		Read:   resourceDataprocJobRead,
		Delete: resourceDataprocJobDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			// Ref: https://cloud.google.com/dataproc/docs/reference/rest/v1/projects.regions.jobs#JobReference
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "global",
				ForceNew: true,
			},

			// If a job is still running, trying to delete a job will fail. Setting
			// this flag to true however will force the deletion by first cancelling
			// the job and then deleting it
			"force_delete": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},

			"reference": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": {
							Type:         schema.TypeString,
							Description:  "The job ID, which must be unique within the project. The job ID is generated by the server upon job submission or provided by the user as a means to perform retries without creating duplicate jobs",
							Optional:     true,
							ForceNew:     true,
							Computed:     true,
							ValidateFunc: validateRegexp("^[a-zA-Z0-9_-]{1,100}$"),
						},
					},
				},
			},

			"placement": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"cluster_name": {
							Type:        schema.TypeString,
							Description: "The name of the cluster where the job will be submitted",
							Required:    true,
							ForceNew:    true,
						},
						"cluster_uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Output-only. A cluster UUID generated by the Cloud Dataproc service when the job is submitted",
						},
					},
				},
			},

			"status": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"state": {
							Type:        schema.TypeString,
							Description: "Output-only. A state message specifying the overall job state",
							Computed:    true,
						},
						"details": {
							Type:        schema.TypeString,
							Description: "Output-only. Optional job state details, such as an error description if the state is ERROR",
							Computed:    true,
						},
						"state_start_time": {
							Type:        schema.TypeString,
							Description: "Output-only. The time when this state was entered",
							Computed:    true,
						},
						"substate": {
							Type:        schema.TypeString,
							Description: "Output-only. Additional state information, which includes status reported by the agent",
							Computed:    true,
						},
					},
				},
			},

			"driver_output_resource_uri": {
				Type:        schema.TypeString,
				Description: "Output-only. A URI pointing to the location of the stdout of the job's driver program",
				Computed:    true,
			},

			"driver_controls_files_uri": {
				Type:        schema.TypeString,
				Description: "Output-only. If present, the location of miscellaneous control files which may be used as part of job setup and handling. If not present, control files may be placed in the same location as driver_output_uri.",
				Computed:    true,
			},

			"labels": {
				Type:        schema.TypeMap,
				Description: "Optional. The labels to associate with this job.",
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"scheduling": {
				Type:        schema.TypeList,
				Description: "Optional. Job scheduling configuration.",
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_failures_per_hour": {
							Type:         schema.TypeInt,
							Description:  "Maximum number of times per hour a driver may be restarted as a result of driver terminating with non-zero code before job is reported failed.",
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntAtMost(10),
						},
					},
				},
			},

			"pyspark_config":  pySparkSchema,
			"spark_config":    sparkSchema,
			"hadoop_config":   hadoopSchema,
			"hive_config":     hiveSchema,
			"pig_config":      pigSchema,
			"sparksql_config": sparkSqlSchema,
		},
	}
}

func resourceDataprocJobUpdate(d *schema.ResourceData, meta interface{}) error {

	// The only updatable value is currently 'force_delete' which is a local
	// only value therefore we don't need to make any GCP calls to update this.

	return resourceDataprocJobRead(d, meta)
}

func resourceDataprocJobCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	jobConfCount := 0
	clusterName := d.Get("placement.0.cluster_name").(string)
	region := d.Get("region").(string)

	submitReq := &dataproc.SubmitJobRequest{
		Job: &dataproc.Job{
			Placement: &dataproc.JobPlacement{
				ClusterName: clusterName,
			},
			Reference: &dataproc.JobReference{
				ProjectId: project,
			},
		},
	}

	if v, ok := d.GetOk("reference.0.job_id"); ok {
		submitReq.Job.Reference.JobId = v.(string)
	}
	if _, ok := d.GetOk("labels"); ok {
		submitReq.Job.Labels = expandLabels(d)
	}

	if v, ok := d.GetOk("pyspark_config"); ok {
		jobConfCount++
		config := extractFirstMapConfig(v.([]interface{}))
		submitReq.Job.PysparkJob = expandPySparkJob(config)
	}

	if v, ok := d.GetOk("spark_config"); ok {
		jobConfCount++
		config := extractFirstMapConfig(v.([]interface{}))
		submitReq.Job.SparkJob = expandSparkJob(config)
	}

	if v, ok := d.GetOk("hadoop_config"); ok {
		jobConfCount++
		config := extractFirstMapConfig(v.([]interface{}))
		submitReq.Job.HadoopJob = expandHadoopJob(config)
	}

	if v, ok := d.GetOk("hive_config"); ok {
		jobConfCount++
		config := extractFirstMapConfig(v.([]interface{}))
		submitReq.Job.HiveJob = expandHiveJob(config)
	}

	if v, ok := d.GetOk("pig_config"); ok {
		jobConfCount++
		config := extractFirstMapConfig(v.([]interface{}))
		submitReq.Job.PigJob = expandPigJob(config)
	}

	if v, ok := d.GetOk("sparksql_config"); ok {
		jobConfCount++
		config := extractFirstMapConfig(v.([]interface{}))
		submitReq.Job.SparkSqlJob = expandSparkSqlJob(config)
	}

	if jobConfCount != 1 {
		return fmt.Errorf("You must define and configure exactly one xxx_config block")
	}

	// Submit the job
	job, err := config.clientDataproc.Projects.Regions.Jobs.Submit(
		project, region, submitReq).Do()
	if err != nil {
		return err
	}
	d.SetId(job.Reference.JobId)

	log.Printf("[INFO] Dataproc job %s has been submitted", job.Reference.JobId)
	return resourceDataprocJobRead(d, meta)
}

func resourceDataprocJobRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := d.Get("region").(string)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	job, err := config.clientDataproc.Projects.Regions.Jobs.Get(
		project, region, d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Dataproc Job %q", d.Id()))
	}

	d.Set("force_delete", d.Get("force_delete"))
	d.Set("labels", job.Labels)
	d.Set("driver_output_resource_uri", job.DriverOutputResourceUri)
	d.Set("driver_controls_files_uri", job.DriverControlFilesUri)

	d.Set("placement", flattenJobPlacement(job.Placement))
	d.Set("status", flattenJobStatus(job.Status))
	d.Set("reference", flattenJobReference(job.Reference))

	if job.PysparkJob != nil {
		d.Set("pyspark_config", flattenPySparkJob(job.PysparkJob))
	}
	if job.SparkJob != nil {
		d.Set("spark_config", flattenSparkJob(job.SparkJob))
	}
	if job.HadoopJob != nil {
		d.Set("hadoop_config", flattenHadoopJob(job.HadoopJob))
	}
	if job.HiveJob != nil {
		d.Set("hive_config", flattenHiveJob(job.HiveJob))
	}
	if job.PigJob != nil {
		d.Set("pig_config", flattenPigJob(job.PigJob))
	}
	if job.SparkSqlJob != nil {
		d.Set("sparksql_config", flattenSparkSqlJob(job.SparkSqlJob))
	}
	return nil
}

func resourceDataprocJobDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region := d.Get("region").(string)
	forceDelete := d.Get("force_delete").(bool)
	timeoutInMinutes := int(d.Timeout(schema.TimeoutDelete).Minutes())

	if forceDelete {
		log.Printf("[DEBUG] Attempting to first cancel Dataproc job %s if it's still running ...", d.Id())

		config.clientDataproc.Projects.Regions.Jobs.Cancel(
			project, region, d.Id(), &dataproc.CancelJobRequest{}).Do()
		// ignore error if we get one - job may be finished already and not need to
		// be cancelled. We do however wait for the state to be one that is
		// at least not active
		waitErr := dataprocJobOperationWait(config, region, project, d.Id(),
			"Cancelling Dataproc job", timeoutInMinutes, 1)
		if waitErr != nil {
			return waitErr
		}

	}

	log.Printf("[DEBUG] Deleting Dataproc job %s", d.Id())
	_, err = config.clientDataproc.Projects.Regions.Jobs.Delete(
		project, region, d.Id()).Do()
	if err != nil {
		return err
	}

	waitErr := dataprocDeleteOperationWait(config, region, project, d.Id(),
		"Deleting Dataproc job", timeoutInMinutes, 1)
	if waitErr != nil {
		return waitErr
	}

	log.Printf("[INFO] Dataproc job %s has been deleted", d.Id())
	d.SetId("")

	return nil
}

// ---- PySpark Job ----

func loggingConfig() *schema.Schema {

	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "The runtime logging config of the job",
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{

				"driver_log_levels": {
					Type:        schema.TypeMap,
					Description: "Optional. The per-package log levels for the driver. This may include 'root' package name to configure rootLogger. Examples: 'com.google = FATAL', 'root = INFO', 'org.apache = DEBUG'.",
					Optional:    true,
					ForceNew:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}

}

var pySparkSchema = &schema.Schema{
	Type:          schema.TypeList,
	Optional:      true,
	ForceNew:      true,
	MaxItems:      1,
	ConflictsWith: []string{"spark_config", "hadoop_config", "hive_config", "pig_config", "sparksql_config"},
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{

			"main_python_file_uri": {
				Type:        schema.TypeString,
				Description: "Required. The HCFS URI of the main Python file to use as the driver. Must be a .py file",
				Required:    true,
				ForceNew:    true,
			},

			"args": {
				Type:        schema.TypeList,
				Description: "Optional. The arguments to pass to the driver. Do not include arguments, such as --conf, that can be set as job properties, since a collision may occur that causes an incorrect job submission",
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"python_file_uris": {
				Type:        schema.TypeList,
				Description: "Optional. HCFS file URIs of Python files to pass to the PySpark framework. Supported file types: .py, .egg, and .zip",
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"jar_file_uris": {
				Type:        schema.TypeList,
				Description: "Optional. HCFS URIs of jar files to add to the CLASSPATHs of the Python driver and tasks",
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"file_uris": {
				Type:        schema.TypeList,
				Description: "Optional. HCFS URIs of files to be copied to the working directory of Python drivers and distributed tasks. Useful for naively parallel tasks",
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"archive_uris": {
				Type:        schema.TypeList,
				Description: "Optional. HCFS URIs of archives to be extracted in the working directory of .jar, .tar, .tar.gz, .tgz, and .zip",
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"properties": {
				Type:        schema.TypeMap,
				Description: "Optional. A mapping of property names to values, used to configure PySpark. Properties that conflict with values set by the Cloud Dataproc API may be overwritten. Can include properties set in /etc/spark/conf/spark-defaults.conf and classes in user code",
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"logging_config": loggingConfig(),
		},
	},
}

func flattenPySparkJob(job *dataproc.PySparkJob) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"main_python_file_uri": job.MainPythonFileUri,
			"args":                 job.Args,
			"python_file_uris":     job.PythonFileUris,
			"jar_file_uris":        job.JarFileUris,
			"file_uris":            job.FileUris,
			"archive_uris":         job.ArchiveUris,
			"properties":           job.Properties,
			"logging_config":       flattenLoggingConfig(job.LoggingConfig),
		},
	}
}

func expandPySparkJob(config map[string]interface{}) *dataproc.PySparkJob {

	job := &dataproc.PySparkJob{}
	if v, ok := config["main_python_file_uri"]; ok {
		job.MainPythonFileUri = v.(string)
	}
	if v, ok := config["args"]; ok {
		job.Args = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["python_file_uris"]; ok {
		job.Args = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["jar_file_uris"]; ok {
		job.Args = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["file_uris"]; ok {
		job.Args = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["archive_uris"]; ok {
		job.Args = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["properties"]; ok {
		job.Properties = convertStringMap(v.(map[string]interface{}))
	}
	if v, ok := config["logging_config"]; ok {
		config := extractFirstMapConfig(v.([]interface{}))
		job.LoggingConfig = expandLoggingConfig(config)
	}

	return job

}

// ---- Spark Job ----

var sparkSchema = &schema.Schema{
	Type:          schema.TypeList,
	Optional:      true,
	ForceNew:      true,
	MaxItems:      1,
	ConflictsWith: []string{"pyspark_config", "hadoop_config", "hive_config", "pig_config", "sparksql_config"},
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{

			// main driver: can be only one of the class | jar_file
			"main_class": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"spark_config.0.main_jar_file_uri"},
			},

			"main_jar_file_uri": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"spark_config.0.main_class"},
			},

			"args": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"jar_file_uris": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"file_uris": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"archive_uris": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"properties": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"logging_config": loggingConfig(),
		},
	},
}

func flattenSparkJob(job *dataproc.SparkJob) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"main_class":        job.MainClass,
			"main_jar_file_uri": job.MainJarFileUri,
			"args":              job.Args,
			"jar_file_uris":     job.JarFileUris,
			"file_uris":         job.FileUris,
			"archive_uris":      job.ArchiveUris,
			"properties":        job.Properties,
			"logging_config":    flattenLoggingConfig(job.LoggingConfig),
		},
	}
}

func expandSparkJob(config map[string]interface{}) *dataproc.SparkJob {

	job := &dataproc.SparkJob{}
	if v, ok := config["main_class"]; ok {
		job.MainClass = v.(string)
	}
	if v, ok := config["main_jar_file_uri"]; ok {
		job.MainJarFileUri = v.(string)
	}

	if v, ok := config["args"]; ok {
		job.Args = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["jar_file_uris"]; ok {
		job.JarFileUris = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["file_uris"]; ok {
		job.FileUris = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["archive_uris"]; ok {
		job.ArchiveUris = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["properties"]; ok {
		job.Properties = convertStringMap(v.(map[string]interface{}))
	}
	if v, ok := config["logging_config"]; ok {
		config := extractFirstMapConfig(v.([]interface{}))
		job.LoggingConfig = expandLoggingConfig(config)
	}

	return job

}

// ---- Hadoop Job ----

var hadoopSchema = &schema.Schema{
	Type:          schema.TypeList,
	Optional:      true,
	ForceNew:      true,
	MaxItems:      1,
	ConflictsWith: []string{"spark_config", "pyspark_config", "hive_config", "pig_config", "sparksql_config"},
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{

			// main driver: can be only one of the main_class | main_jar_file_uri
			"main_class": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"hadoop_config.0.main_jar_file_uri"},
			},

			"main_jar_file_uri": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"hadoop_config.0.main_class"},
			},

			"args": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"jar_file_uris": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"file_uris": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"archive_uris": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"properties": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"logging_config": loggingConfig(),
		},
	},
}

func flattenHadoopJob(job *dataproc.HadoopJob) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"main_class":        job.MainClass,
			"main_jar_file_uri": job.MainJarFileUri,
			"args":              job.Args,
			"jar_file_uris":     job.JarFileUris,
			"file_uris":         job.FileUris,
			"archive_uris":      job.ArchiveUris,
			"properties":        job.Properties,
			"logging_config":    flattenLoggingConfig(job.LoggingConfig),
		},
	}
}

func expandHadoopJob(config map[string]interface{}) *dataproc.HadoopJob {

	job := &dataproc.HadoopJob{}
	if v, ok := config["main_class"]; ok {
		job.MainClass = v.(string)
	}
	if v, ok := config["main_jar_file_uri"]; ok {
		job.MainJarFileUri = v.(string)
	}

	if v, ok := config["args"]; ok {
		job.Args = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["jar_file_uris"]; ok {
		job.JarFileUris = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["file_uris"]; ok {
		job.FileUris = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["archive_uris"]; ok {
		job.ArchiveUris = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["properties"]; ok {
		job.Properties = convertStringMap(v.(map[string]interface{}))
	}
	if v, ok := config["logging_config"]; ok {
		config := extractFirstMapConfig(v.([]interface{}))
		job.LoggingConfig = expandLoggingConfig(config)
	}

	return job

}

// ---- Hive Job ----

var hiveSchema = &schema.Schema{
	Type:          schema.TypeList,
	Optional:      true,
	ForceNew:      true,
	MaxItems:      1,
	ConflictsWith: []string{"spark_config", "pyspark_config", "hadoop_config", "pig_config", "sparksql_config"},
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{

			// main query: can be only one of query_list | query_file_uri
			"query_list": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"hive_config.0.query_file_uri"},
			},

			"query_file_uri": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"hive_config.0.query_list"},
			},

			"continue_on_failure": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"script_variables": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"properties": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"jar_file_uris": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	},
}

func flattenHiveJob(job *dataproc.HiveJob) []map[string]interface{} {

	queries := []string{}
	if job.QueryList != nil {
		queries = job.QueryList.Queries
	}
	return []map[string]interface{}{
		{
			"query_list":          queries,
			"query_file_uri":      job.QueryFileUri,
			"continue_on_failure": job.ContinueOnFailure,
			"script_variables":    job.ScriptVariables,
			"properties":          job.Properties,
			"jar_file_uris":       job.JarFileUris,
		},
	}
}

func expandHiveJob(config map[string]interface{}) *dataproc.HiveJob {

	job := &dataproc.HiveJob{}
	if v, ok := config["query_file_uri"]; ok {
		job.QueryFileUri = v.(string)
	}
	if v, ok := config["query_list"]; ok {
		job.QueryList = &dataproc.QueryList{
			Queries: convertStringArr(v.([]interface{})),
		}
	}
	if v, ok := config["continue_on_failure"]; ok {
		job.ContinueOnFailure = v.(bool)
	}
	if v, ok := config["script_variables"]; ok {
		job.ScriptVariables = convertStringMap(v.(map[string]interface{}))
	}
	if v, ok := config["jar_file_uris"]; ok {
		job.JarFileUris = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["properties"]; ok {
		job.Properties = convertStringMap(v.(map[string]interface{}))
	}

	return job

}

// ---- Pig Job ----

var pigSchema = &schema.Schema{
	Type:          schema.TypeList,
	Optional:      true,
	ForceNew:      true,
	MaxItems:      1,
	ConflictsWith: []string{"spark_config", "pyspark_config", "hadoop_config", "hive_config", "sparksql_config"},
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{

			// main query: can be only one of query_list | query_file_uri
			"query_list": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"pig_config.0.query_file_uri"},
			},

			"query_file_uri": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"pig_config.0.query_list"},
			},

			"continue_on_failure": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"script_variables": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"properties": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"jar_file_uris": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"logging_config": loggingConfig(),
		},
	},
}

func flattenPigJob(job *dataproc.PigJob) []map[string]interface{} {

	queries := []string{}
	if job.QueryList != nil {
		queries = job.QueryList.Queries
	}
	return []map[string]interface{}{
		{
			"query_list":          queries,
			"query_file_uri":      job.QueryFileUri,
			"continue_on_failure": job.ContinueOnFailure,
			"script_variables":    job.ScriptVariables,
			"properties":          job.Properties,
			"jar_file_uris":       job.JarFileUris,
		},
	}
}

func expandPigJob(config map[string]interface{}) *dataproc.PigJob {

	job := &dataproc.PigJob{}
	if v, ok := config["query_file_uri"]; ok {
		job.QueryFileUri = v.(string)
	}
	if v, ok := config["query_list"]; ok {
		job.QueryList = &dataproc.QueryList{
			Queries: convertStringArr(v.([]interface{})),
		}
	}
	if v, ok := config["continue_on_failure"]; ok {
		job.ContinueOnFailure = v.(bool)
	}
	if v, ok := config["script_variables"]; ok {
		job.ScriptVariables = convertStringMap(v.(map[string]interface{}))
	}
	if v, ok := config["jar_file_uris"]; ok {
		job.JarFileUris = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["properties"]; ok {
		job.Properties = convertStringMap(v.(map[string]interface{}))
	}

	return job

}

// ---- Spark SQL Job ----

var sparkSqlSchema = &schema.Schema{
	Type:          schema.TypeList,
	Optional:      true,
	ForceNew:      true,
	MaxItems:      1,
	ConflictsWith: []string{"spark_config", "pyspark_config", "hadoop_config", "hive_config", "pig_config"},
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{

			// main query: can be only one of query_list | query_file_uri
			"query_list": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"pig_config.0.query_file_uri"},
			},

			"query_file_uri": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"pig_config.0.query_list"},
			},

			"script_variables": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"properties": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"jar_file_uris": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"logging_config": loggingConfig(),
		},
	},
}

func flattenSparkSqlJob(job *dataproc.SparkSqlJob) []map[string]interface{} {

	queries := []string{}
	if job.QueryList != nil {
		queries = job.QueryList.Queries
	}
	return []map[string]interface{}{
		{
			"query_list":       queries,
			"query_file_uri":   job.QueryFileUri,
			"script_variables": job.ScriptVariables,
			"properties":       job.Properties,
			"jar_file_uris":    job.JarFileUris,
		},
	}
}

func expandSparkSqlJob(config map[string]interface{}) *dataproc.SparkSqlJob {

	job := &dataproc.SparkSqlJob{}
	if v, ok := config["query_file_uri"]; ok {
		job.QueryFileUri = v.(string)
	}
	if v, ok := config["query_list"]; ok {
		job.QueryList = &dataproc.QueryList{
			Queries: convertStringArr(v.([]interface{})),
		}
	}
	if v, ok := config["script_variables"]; ok {
		job.ScriptVariables = convertStringMap(v.(map[string]interface{}))
	}
	if v, ok := config["jar_file_uris"]; ok {
		job.JarFileUris = convertStringArr(v.([]interface{}))
	}
	if v, ok := config["properties"]; ok {
		job.Properties = convertStringMap(v.(map[string]interface{}))
	}

	return job

}

// ---- Other flatten / expand methods ----

func expandLoggingConfig(config map[string]interface{}) *dataproc.LoggingConfig {
	conf := &dataproc.LoggingConfig{}
	if v, ok := config["driver_log_levels"]; ok {
		conf.DriverLogLevels = convertStringMap(v.(map[string]interface{}))
	}
	return conf
}

func flattenLoggingConfig(l *dataproc.LoggingConfig) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"driver_log_levels": l.DriverLogLevels,
		},
	}
}

func flattenJobReference(r *dataproc.JobReference) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"job_id": r.JobId,
		},
	}
}

func flattenJobStatus(s *dataproc.JobStatus) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"state":            s.State,
			"details":          s.Details,
			"state_start_time": s.StateStartTime,
			"substate":         s.Substate,
		},
	}
}

func flattenJobPlacement(jp *dataproc.JobPlacement) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"cluster_name": jp.ClusterName,
			"cluster_uuid": jp.ClusterUuid,
		},
	}
}
