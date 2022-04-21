
//export const URL = "https://backend-sergioandresestrada.cloud.okteto.net"
export const URL = "http://192.168.1.208:12345"

export function isValidFile(file: File) : boolean{
    var acceptedTypes = ["pdf", "stl"]

    var re = /(?:\.([^.]+))?$/

    var result = re.exec(file.name)
    if (result === null) return false

    var fileExtension = result[1]
    if (acceptedTypes.indexOf(fileExtension) === -1){
        return false
    }

    return true

}

const REGEX_IPAddress = /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/

export function validateIP(ip: string | undefined) : boolean{
    if (ip == null) return false
    if(REGEX_IPAddress.test(ip)){
        return true
    }
    return false
}

// Just delete the type of information prefix and the file extension
export function beautifyFileName(fn : string) : string {
    return fn.replace("Jobs-", "").replace("Identification-","").replace(".json", "")
}
 