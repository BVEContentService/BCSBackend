## OpenBVE MCS REST API Documentation
This documentation provides information on making requests to the REST API server,
as well as parsing the response returned from it.
## Conventions
1. **API Root Path URL**:  `https://api.bvecs.tk/v1`
2. Most Endpoints support both XML and JSON request body. **Content-Type** header 
   of `application/xml` and `application/json` are both supported.
   XML is assumed if **Content-Type** is not specified.
3. Most Endpoints support both XML and JSON response. **Accept** header
   of `application/xml` and `application/json` are both supported.
   XML is used if **Accept** is not specified.
4. JWT Endpoints are currently **JSON-only** for all requests and responses.
## Models
### Package Structure
|Field Name |Type      | Description                                           |
|-----------|----------|-------------------------------------------------------|
|GUID       |str       | The Globally Unique Identifier of the package.        |
|Identifier |str       | The Human-readable unique name of the package.        |
|Name       |String3   | The display name of the package.                      |
|UploaderID |uint      | The UID of the uploader.                              |
|_Uploader_ |Uploader? | The detailed information of the uploader.             |
|Author     |Developer?| About original author, if the package is not original |
|Homepage   |str?      | A URL containing more info about the package.         |
|Thumbnail  |str?      | A URL providing a thumbnail image of the package.     |
|ThumbnailLQ|str?      | Same as above, but with a tinier file size.           |
|Description|str?      | A HTML document containing more description.          |
|_Files_    |File[]?   | A list of files attached to this package.             |
> Question Mark: The field might not be included in response if not applicable.  
> _Italic Line_: The field is not included in the list of all packages.

### API Error Structure and Codes
An Error Structure will be returned as response when an API error arose.

|FieldName|Type| Description                                          |
|---------|----|------------------------------------------------------|
|ErrorCode|int | The identifier of this error, as described below.    |
|Msg      |str | A message providing extra description of this error. |
|Request  |str | The method and path requested when the error arose.  |

|ErrorCode|HTTP| Description                                          |
|---------|----|------------------------------------------------------|
|   101   |404 | The path requested is not an API Endpoint.           |
|   102   |405 | This Endpoint does not support the requested Method. |
|   111   |400 | The Request Body is not in the required format.      |
|   112   |422 | Some fields in the Request Body failed validation.   |
|   201   |404 | The Resource requested does not exist.               |
|   202   |422 | The Identifier Name is already taken.                |
|   203   |422 | The Identifier GUID is already taken.                |
|   204   |422 | Duplicate Files with same Version on same Platform.  |
|   211   |401 | Generic JWT error. More info is in the message.      |
|   212   |403 | Incorrect Username or Password at login.             |
|   221   |403 | Such an operation is beyond your authority.          |
| No Body |500 | Server crashed. Contact the developer immediately.   |
> Error #211 is most commonly caused by token expiration.
> Token renewal or another login can be attempted.

## Endpoints

## File Validation

  **Validated** is a property that can be modified by a moderator, to indicate 
  whether a File shall be displayed to the public.  
  However, not all Files uploaded to the platform needs to be validated. The
  genres of Files that needs to be validated can be specified in the server 
  configuration.
  Therefore, if a File is not required to be validated (**NeedValidation**=False),
  then it can be displayed regardless of the value of the **Validated** field.
  
  Here describes the handling of the two aforementioned fields on different API
  Endpoints. Please notice the difference:
  
  - GET /packages (?platform Filter) :  
    Only validated files will be considered when collecting the available
    platforms of a package. Files that do not need validation are considered
    validated.
    
  - GET /files (?validated Filter) :  
    Files that do not need validation are considered validated. **Platform** field
    is not affected.
    
  - GET /packages/:id :  
    **Platforms** Field: Only validated files will be considered when collecting
    the available platforms of a package. Files that do not need validation are 
    considered validated.  
    **Files[]** Array: All packages regardless of whether it is validated are all 
    included. Both fields are included in the Files[] response and the
    **Validated** field is not changed to accommodate for the **NeedValidation**
    field. It is up to the client to accommodate for both the fields and decide 
    if the unvalidated packages should be hidden.