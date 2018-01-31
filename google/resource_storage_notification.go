package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"google.golang.org/api/storage/v1"
)

func resourceStorageNotification() *schema.Resource {
	return &schema.Resource{
		Create: resourceStorageNotificationCreate,
		Read:   resourceStorageNotificationRead,
		Delete: resourceStorageNotificationDelete,
		Importer: &schema.ResourceImporter{
			State: resourceStorageNotificationImportState,
		},

		Schema: map[string]*schema.Schema{
			"bucket": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"payload_format": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"JSON_API_V1", "NONE"}, false),
			},

			"topic": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"custom_attributes": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"event_types": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"OBJECT_FINALIZE", "OBJECT_METADATA_UPDATE", "OBJECT_DELETE", "OBJECT_ARCHIVE"},
						false),
				},
			},

			"prefix": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceStorageNotificationCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	bucket := d.Get("bucket").(string)

	storageNotification := &storage.Notification{
		CustomAttributes: expandStringMap(d, "custom_attributes"),
		EventTypes:       convertStringSet(d.Get("event_types").(*schema.Set)),
		ObjectNamePrefix: d.Get("prefix").(string),
		PayloadFormat:    d.Get("payload_format").(string),
		Topic:            d.Get("topic").(string),
	}

	res, err := config.clientStorage.Notifications.Insert(bucket, storageNotification).Do()

	if err != nil {
		return fmt.Errorf("Error creating notification config for bucket %s: %v", bucket, err)
	}

	d.SetId(fmt.Sprintf("%s/notificationConfigs/%s", bucket, res.Id))

	return resourceStorageNotificationRead(d, meta)
}

func resourceStorageNotificationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	bucket := d.Get("bucket").(string)
	notificationID := resourceStorageNotificationParseID(d.Id())

	res, err := config.clientStorage.Notifications.Get(bucket, notificationID).Do()

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Notification configuration %s for bucket %s", notificationID, bucket))
	}

	d.Set("payload_format", res.PayloadFormat)
	d.Set("topic", strings.SplitAfter(res.Topic, "//pubsub.googleapis.com/")[1])
	d.Set("prefix", res.ObjectNamePrefix)
	d.Set("event_types", res.EventTypes)
	d.Set("self_link", res.SelfLink)
	d.Set("custom_attributes", res.CustomAttributes)

	return nil
}

func resourceStorageNotificationDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	bucket := d.Get("bucket").(string)

	notificationID := resourceStorageNotificationParseID(d.Id())

	err := config.clientStorage.Notifications.Delete(bucket, notificationID).Do()

	if err != nil {
		return fmt.Errorf("Error deleting notification configuration %s for bucket %s: %v", notificationID, bucket, err)
	}

	return nil
}

func resourceStorageNotificationImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 || parts[1] != "notificationConfigs" {
		return nil, fmt.Errorf("Invalid notification configuration specifier. Expecting {bucket}/notificationConfigs/{id}")
	}

	d.Set("bucket", parts[0])

	if err := resourceStorageNotificationRead(d, meta); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func resourceStorageNotificationParseID(id string) string {
	//NotificationID
	parts := strings.Split(id, "/")
	return parts[2]
}
