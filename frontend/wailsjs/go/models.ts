export namespace main {
	
	export class Config {
	    appId: string;
	    appSecret: string;
	    appKey: string;
	    clusterId: string;
	    supaBaseApiKey: string;
	    supaBaseUrl: string;
	    rabbitMqAdmin: string;
	    rabbitMqPassword: string;
	    rabbitMqHost: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.appId = source["appId"];
	        this.appSecret = source["appSecret"];
	        this.appKey = source["appKey"];
	        this.clusterId = source["clusterId"];
	        this.supaBaseApiKey = source["supaBaseApiKey"];
	        this.supaBaseUrl = source["supaBaseUrl"];
	        this.rabbitMqAdmin = source["rabbitMqAdmin"];
	        this.rabbitMqPassword = source["rabbitMqPassword"];
	        this.rabbitMqHost = source["rabbitMqHost"];
	    }
	}

}

