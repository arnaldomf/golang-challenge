# Protocol
## Version: 0.808-alpha

`HEADER`: **SPLICE**|0|0|0|0|0|0|0|*Bytes until the end*|

`VERSION`: <**version-string**>

`TRACK`: |*ID*|0|0|0|*Sound Name Size*|<**Sound Name**>|x|x|x|x|x|x|x|x|x|x|x|x|x|x|x|x|

###Specs:
1. Header has 14 bytes
2. Version has a variable length
3. Each track starts with its ID and has 16 **steps**
4. The sound name size is variable
