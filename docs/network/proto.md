# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [network/networkservice.proto](#network_networkservice-proto)
    - [InitialConfiguration](#networkservice-InitialConfiguration)
    - [NetworkAttachment](#networkservice-NetworkAttachment)
    - [NetworkAttachmentName](#networkservice-NetworkAttachmentName)
    - [NetworkAttachments](#networkservice-NetworkAttachments)
    - [NicSelector](#networkservice-NicSelector)
    - [Resource](#networkservice-Resource)
    - [ResourceConfig](#networkservice-ResourceConfig)
    - [ResourceConfigs](#networkservice-ResourceConfigs)
    - [ResourceName](#networkservice-ResourceName)
    - [ResourceSpec](#networkservice-ResourceSpec)
    - [ResourceStatus](#networkservice-ResourceStatus)
    - [Resources](#networkservice-Resources)
    - [VFResourceStatus](#networkservice-VFResourceStatus)
  
    - [NetworkAttachmentService](#networkservice-NetworkAttachmentService)
    - [ResourceService](#networkservice-ResourceService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="network_networkservice-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## network/networkservice.proto



<a name="networkservice-InitialConfiguration"></a>

### InitialConfiguration



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resourceConfigs | [ResourceConfig](#networkservice-ResourceConfig) | repeated | list of resource configurations |
| networkattachments | [NetworkAttachment](#networkservice-NetworkAttachment) | repeated | list of network attachments |






<a name="networkservice-NetworkAttachment"></a>

### NetworkAttachment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| resourceName | [string](#string) |  |  |
| config | [google.protobuf.Struct](#google-protobuf-Struct) |  |  |






<a name="networkservice-NetworkAttachmentName"></a>

### NetworkAttachmentName



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |






<a name="networkservice-NetworkAttachments"></a>

### NetworkAttachments



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| networkattachments | [NetworkAttachment](#networkservice-NetworkAttachment) | repeated |  |






<a name="networkservice-NicSelector"></a>

### NicSelector



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vendors | [string](#string) | repeated |  |
| drivers | [string](#string) | repeated |  |
| devices | [string](#string) | repeated |  |
| pfNames | [string](#string) | repeated |  |






<a name="networkservice-Resource"></a>

### Resource



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| spec | [ResourceSpec](#networkservice-ResourceSpec) |  | resource configuration spec |
| status | [ResourceStatus](#networkservice-ResourceStatus) | repeated | discovered status corresponding to the spec |






<a name="networkservice-ResourceConfig"></a>

### ResourceConfig



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| mtu | [uint32](#uint32) |  |  |
| numVfs | [uint32](#uint32) |  |  |
| needVhostNet | [bool](#bool) |  |  |
| nicSelector | [NicSelector](#networkservice-NicSelector) |  |  |
| deviceType | [string](#string) |  |  |






<a name="networkservice-ResourceConfigs"></a>

### ResourceConfigs



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resourceConfigs | [ResourceConfig](#networkservice-ResourceConfig) | repeated |  |






<a name="networkservice-ResourceName"></a>

### ResourceName



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |






<a name="networkservice-ResourceSpec"></a>

### ResourceSpec



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| mtu | [uint32](#uint32) |  |  |
| numVfs | [uint32](#uint32) |  |  |
| needVhostNet | [bool](#bool) |  |  |
| devices | [string](#string) | repeated |  |






<a name="networkservice-ResourceStatus"></a>

### ResourceStatus



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| mtu | [uint32](#uint32) |  |  |
| numVfs | [uint32](#uint32) |  |  |
| mac | [string](#string) |  |  |
| vendor | [string](#string) |  |  |
| driver | [string](#string) |  |  |
| device | [string](#string) |  |  |
| vfs | [VFResourceStatus](#networkservice-VFResourceStatus) | repeated |  |






<a name="networkservice-Resources"></a>

### Resources



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resources | [Resource](#networkservice-Resource) | repeated |  |






<a name="networkservice-VFResourceStatus"></a>

### VFResourceStatus



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| mac | [string](#string) |  |  |
| vendor | [string](#string) |  |  |
| driver | [string](#string) |  |  |
| device | [string](#string) |  |  |





 

 

 


<a name="networkservice-NetworkAttachmentService"></a>

### NetworkAttachmentService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateNetworkAttachment | [NetworkAttachment](#networkservice-NetworkAttachment) | [.google.protobuf.Empty](#google-protobuf-Empty) |  |
| DeleteNetworkAttachment | [NetworkAttachmentName](#networkservice-NetworkAttachmentName) | [.google.protobuf.Empty](#google-protobuf-Empty) |  |
| GetAllNetworkAttachments | [.google.protobuf.Empty](#google-protobuf-Empty) | [NetworkAttachments](#networkservice-NetworkAttachments) |  |
| GetNetworkAttachment | [NetworkAttachmentName](#networkservice-NetworkAttachmentName) | [NetworkAttachment](#networkservice-NetworkAttachment) |  |


<a name="networkservice-ResourceService"></a>

### ResourceService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateResourceConfig | [ResourceConfig](#networkservice-ResourceConfig) | [Resource](#networkservice-Resource) |  |
| DeleteResourceConfig | [ResourceName](#networkservice-ResourceName) | [.google.protobuf.Empty](#google-protobuf-Empty) |  |
| GetAllResourceConfigs | [.google.protobuf.Empty](#google-protobuf-Empty) | [ResourceConfigs](#networkservice-ResourceConfigs) |  |
| GetResourceConfig | [ResourceName](#networkservice-ResourceName) | [ResourceConfig](#networkservice-ResourceConfig) |  |
| GetAllResources | [.google.protobuf.Empty](#google-protobuf-Empty) | [Resources](#networkservice-Resources) |  |
| GetResource | [ResourceName](#networkservice-ResourceName) | [Resource](#networkservice-Resource) |  |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

