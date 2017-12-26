package google

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/cloudfunctions/v1"
	//"github.com/hashicorp/terraform/helper/resource"
	//"bytes"
)

const (
	CLOUDFUNCTIONS_FULL_NAME   = 0
	CLOUDFUNCTIONS_REGION_ONLY = 1
)

//Function would return formatted string to be used in API calls for name or location
// Arguments:
//  funcType: {CLOUDFUNCTIONS_REGION_ONLY, CLOUDFUNCTIONS_FULL_NAME}
//            If specifying CLOUDFUNCTIONS_FULL_NAME string would be in format 'projects/YOUR_PROJECT/locations/REGION/functions/FUNCTION_NAME'
//            If specifying CLOUDFUNCTIONS_REGION_ONLY string would be in format 'projects/YOUR_PROJECT/locations/REGION'
//  projectName: Name of project in Google Cloud
//  region: In which region function is/should be located. NOTE: Not all regions are supported in 2017
//  funcName: name of function. In case CLOUDFUNCTIONS_REGION_ONLY might empty
// Returns:
//  path to function
func createCloudFunctionsPathString(funcType int, projectName string, region string, funcName string) (path string) {
	path = fmt.Sprintf("projects/%s/locations/%s", projectName, region)
	if funcType == CLOUDFUNCTIONS_FULL_NAME {
		path = fmt.Sprintf("%s/functions/%s", path, funcName)
	}
	return
}

//Function would get result for long running operations
//Function might timeout from loop if no result returned in 40 seconds.
//Arguments:
//  opName: string which you get from Create,Delete or other operations to check for result
//  service: Pointer to cloudfunctions API service to use
//Returns:
//  res: result.Done from Operations Result status
//  error: error if any occurred
func getCloudFunctionsOperationsResults(opName string, service *cloudfunctions.Service) (bool, error) {
	var err error

	attempt := 0
	res, err := service.Operations.Get(opName).Do()
	for res.Done == false && attempt < 20 {
		attempt++
		time.Sleep(2000 * time.Millisecond)
		res, err = service.Operations.Get(opName).Do()
		if err != nil {
			return res.Done, err
		}
		if res.Error != nil {
			log.Print(res.Error.Details)
		}
	}
	if res.Done == false && attempt == 20 {
		log.Println("getOperationsResults is not done in 20 attempts...")
	}

	return res.Done, err
}

//Function would extract short function name from long used in Google cloud
//Arguments:
//  fullPath: function name in Google cloud
//Return:
//  functionName: string short function name
//  error: Error if fullPath is not correct Google cloud function name
func getCloudFunctionName(fullPath string) (string, error) {
	allParts, err := splitCloudFunctionFullPath(fullPath)
	if err != nil {
		return "", err
	}
	return allParts[3], nil
}

//Function would extract region from long used in Google cloud
//Arguments:
//  fullPath: function name in Google cloud
//Return:
//  region: string zone in which function is deployed
//  error: Error if fullPath is not correct Google cloud function name
func getCloudFunctionRegion(fullPath string) (string, error) {
	allParts, err := splitCloudFunctionFullPath(fullPath)
	if err != nil {
		return "", err
	}
	return allParts[2], nil
}

//Function would extract project from long used in Google cloud
//Arguments:
//  fullPath: function name in Google cloud
//Return:
//  project: string project in which function is deployed
//  error: Error if fullPath is not correct Google cloud function name
func getCloudFunctionProject(fullPath string) (string, error) {
	allParts, err := splitCloudFunctionFullPath(fullPath)
	if err != nil {
		return "", err
	}
	return allParts[1], nil
}

//Function would split full CloudFunction Path into array of fullName,project,region,funcName
//Arguments:
//  fullPath: function name in Google cloud
//Return:
//  parts: array of strings which would include fullName,project,region,funcName
//  error: Error if fullPath is not correct Google cloud function name
func splitCloudFunctionFullPath(fullPath string) ([]string, error) {
	namePattern := regexp.MustCompile("^projects/([^/]+)/locations/([^/]+)/functions/([^/]+)$")
	if !namePattern.MatchString(fullPath) {
		return nil, fmt.Errorf("%s is not valid CloudFunction full name", fullPath)
	}
	return namePattern.FindStringSubmatch(fullPath), nil
}

//Function would read timeout value from GCloud
//Arguments:
//  timeout: Timeout in string
//Return:
//  timeout int
func readTimeout(s string) (int, error) {
	sRemoved := strings.Replace(s, "s", "", -1)
	return strconv.Atoi(sRemoved)
}
