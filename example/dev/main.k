# The configuration codes in perspective of developers. 
import kam.v1.app_configuration as ac
import service
import service.container as c
import objectstorage.v1.objectstorage

objectstorage: ac.AppConfiguration {
    # Declare the workload configurations. 
    workload: service.Service {
        containers: {
            objectstorage: c.Container {
                image: "hashicorp/http-echo"
                env: {
                    "ECHO_TEXT": "$(KUSION_AWS_S3_BUCKET_DOMAIN_NAME)",
                    "ECHO_TEXT2": "$(KUSION_AWS_S3_BUCKET_REGIONAL_DOMAIN_NAME)"
                }
            }
        }
        replicas: 1
    }
    # Declare the objectstorage module configurations. 
    accessories: {
        "objectstorage": objectstorage.ObjectStorage {
            type: "cloud"
            bucket: "anhtvwordpress"
        }
    }
}
