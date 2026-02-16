export namespace main {
	
	export class CalibratedCalibrationPoint {
	    Id: string;
	    Pan: number;
	    Tilt: number;
	
	    static createFrom(source: any = {}) {
	        return new CalibratedCalibrationPoint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Pan = source["Pan"];
	        this.Tilt = source["Tilt"];
	    }
	}
	export class CalibrationPoint {
	    Id: string;
	    Name: string;
	    X: number;
	    Y: number;
	
	    static createFrom(source: any = {}) {
	        return new CalibrationPoint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Name = source["Name"];
	        this.X = source["X"];
	        this.Y = source["Y"];
	    }
	}
	export class Fixture {
	    Id: string;
	    Name: string;
	    Universe: number;
	    PanAddress: number;
	    FinePanAddress: number;
	    TiltAddress: number;
	    FineTiltAddress: number;
	    Calibration: Record<string, CalibratedCalibrationPoint>;
	
	    static createFrom(source: any = {}) {
	        return new Fixture(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Name = source["Name"];
	        this.Universe = source["Universe"];
	        this.PanAddress = source["PanAddress"];
	        this.FinePanAddress = source["FinePanAddress"];
	        this.TiltAddress = source["TiltAddress"];
	        this.FineTiltAddress = source["FineTiltAddress"];
	        this.Calibration = this.convertValues(source["Calibration"], CalibratedCalibrationPoint, true);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LastSessionInfo {
	    hasLastSession: boolean;
	    configPath: string;
	    configName: string;
	    ipAddress: string;
	    ipAddressValid: boolean;
	    videoSourceId: string;
	    videoSourceLabel: string;
	
	    static createFrom(source: any = {}) {
	        return new LastSessionInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hasLastSession = source["hasLastSession"];
	        this.configPath = source["configPath"];
	        this.configName = source["configName"];
	        this.ipAddress = source["ipAddress"];
	        this.ipAddressValid = source["ipAddressValid"];
	        this.videoSourceId = source["videoSourceId"];
	        this.videoSourceLabel = source["videoSourceLabel"];
	    }
	}
	export class Point {
	    X: number;
	    Y: number;
	
	    static createFrom(source: any = {}) {
	        return new Point(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.X = source["X"];
	        this.Y = source["Y"];
	    }
	}
	export class SACNConfig {
	    IpAddress: string;
	    PossibleIpAddresses: string[];
	    Fps: number;
	    Multicast: boolean;
	    Destinations: string[];
	
	    static createFrom(source: any = {}) {
	        return new SACNConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.IpAddress = source["IpAddress"];
	        this.PossibleIpAddresses = source["PossibleIpAddresses"];
	        this.Fps = source["Fps"];
	        this.Multicast = source["Multicast"];
	        this.Destinations = source["Destinations"];
	    }
	}

}

