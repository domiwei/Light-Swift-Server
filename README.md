TODO:
* Flush dirty data
	- Add dirty flag.
	- Trigger flush when dirty usage exceeds threshold.
	- Traverse all mapping and choose dirty object to write.
* Read data from disk when server up
	- Read all of them.
* Read server data when server up
	- Need write server meta.

* Add LRU-MRU structure to object mapping
* Page out LRU data if mem usage is full
* Cache disk data in memory map.

