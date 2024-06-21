schema ObjectStorage: 
    """ ObjectStorage

    Attributes
    ----------
    type: "local" | "cloud", defaults to Undefined, required. 
        Type defines whether the objectstorage is deployed locally or provided by 
        cloud vendor. 
    bucket: str, defaults to Undefined, required. 
        Version defines the bucket name to use. 


    Examples
    --------
    import objectstorage.v1.objectstorage as os

    accessories: {
        "vietanhtwdk/objectstorage@v0.2.0": os.ObjectStorage {
            type: "cloud"
            bucket: "anhtvwordpress"
        }
    }
    """
    # The deployment mode of the mysql database. 
    type:       "local" | "cloud"

    # The bucket name to use. 
    bucket:    str