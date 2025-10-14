/*
A simple utility package for providing some common cache implementations
Example:

	c := NewRedisCache(RedisConnectionSettings{
		Host: "localhost",
		Port: 6379,
		Password: ""
	})

	c.Set("key1", "test", 5)
	stop := false
	for  stop == false{
		data , err := c.Get("key1")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(d)
	}
*/
package cache
