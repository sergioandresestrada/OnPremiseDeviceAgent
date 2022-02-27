
import React, { FormEvent } from 'react';
import { Form as FormRS, FormGroup, FormText, Input, Label, Modal, ModalBody, ModalHeader, Spinner, Button, ModalFooter} from 'reactstrap';

enum Type {
    HEARTBEAT = "HEARTBEAT",
    JOB = "JOB",
    PLACEHOLDER1 = "PLACEHOLDER1"
}

enum Material {
    "HR PA 11" = "HR PA 11", 
    "HR PA 12GB" = "HR PA 12GB",
    "HR TPA" = "HR TPA", 
    "HR PP" = "HR PP", 
    "HR PA 12" = "HR PA 12"
}

const initialState = {
    message : '',
    type : Type.HEARTBEAT,
    file : undefined,
    material : Material['HR PA 11'],
    IPAddress : "",

    processingJob : false,
    submitOutcome : ''
}

const URL = "https://backend-sergioandresestrada.cloud.okteto.net"
//const URL = "http://192.168.1.208:12345"

const REGEX_IPAddress = /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/

interface IJob{
    message : string,
    type : Type,
    file? : File,
    material? : string
    IPAddress? : string

    processingJob : boolean
    submitOutcome : string
}

class Form extends React.Component<{}, IJob>{ 
    constructor(props : any){
        super(props);
        this.state = initialState
        this.handleSubmit = this.handleSubmit.bind(this)
        this.handleChangeMessage = this.handleChangeMessage.bind(this)
        this.handleChangeType = this.handleChangeType.bind(this)
        this.handleChangeFile = this.handleChangeFile.bind(this)
        this.handleChangeMaterial = this.handleChangeMaterial.bind(this)
        this.handleChangeIP = this.handleChangeIP.bind(this)
        this.validateIP = this.validateIP.bind(this)
    }

    handleSubmit(event : FormEvent){

        event.preventDefault()

        let fetchOptions : object = {}
        let fullURL : string = ""

        switch (this.state.type){
            case "HEARTBEAT":
                fullURL = URL + "/heartbeat"
                fetchOptions = {
                    method: "POST",
                    headers: {
                        "Content-Type": "aplication/json"
                    },
                    body: JSON.stringify({
                        type: this.state.type,
                        message: this.state.message
                    })
                }
                break
                
            case "JOB":
                if (this.state.file === undefined){
                    alert("Error getting the file")
                    return
                }
                if(!this.isValidFile(this.state.file)){
                    alert("Invalid file selected. Please select a PDF or STL file and try again")
                    return
                }
                fullURL = URL + "/job"
                let formData = new FormData()
                let data = JSON.stringify({
                    type: this.state.type,
                    message: this.state.message,
                    material: this.state.material,
                    IPAddress: this.state.IPAddress
                })
                formData.append("data", data)
                formData.append("file", this.state.file)

                fetchOptions = {
                    method: "POST",
                    body: formData
                }
                break
        }

        if (fullURL === "" || fetchOptions === {}) return

        this.setState({
            processingJob : true
        })

        fetch(fullURL, fetchOptions)
        .then(response => {
            let outcome = ""
            switch(response.status){
                case 200:
                    outcome = "New " + this.state.type.charAt(0)+this.state.type.substring(1).toLowerCase()+ " was sent successfully."
                    break
                case 400:
                    outcome = "Bad request. Check the fields and try again."
                    break
                case 500:
                    outcome = "Server error. Try again later."
                    break
            }
            this.setState({
                submitOutcome : outcome,
                processingJob : false
            })
        })
        .catch(error => {
            let outcome = "There was an error connecting to the server, please try again later."
            this.setState({
                submitOutcome : outcome,
                processingJob : false
            })
        })

    }

    resetForm = () => {
        this.setState(initialState)
    }

    isValidFile(file: File) : boolean{
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

    handleChangeMessage(event: React.ChangeEvent<HTMLInputElement>){
        this.setState({
            message : event.target.value
        });
    }

    handleChangeType(event : React.ChangeEvent<HTMLInputElement>){
        this.setState({
            type : event.target.value as Type
        })
    }

    handleChangeFile(event : React.ChangeEvent<HTMLInputElement>){
        if (event.target.files != null){
            this.setState({
                file : event.target.files[0]
            })
        }
    }

    handleChangeMaterial(event: React.ChangeEvent<HTMLInputElement>){
        this.setState({
            material : event.target.value
        })
    }

    handleChangeIP(event: React.ChangeEvent<HTMLInputElement>){
        this.setState({
            IPAddress : event.target.value
        })
    }

    validateIP() : boolean{
        if (this.state.IPAddress == null) return false
        if(REGEX_IPAddress.test(this.state.IPAddress)){
            return true
        }
        return false
    }

    render() {
        return(
            <div>
                <FormRS onSubmit={this.handleSubmit}>
                    <FormGroup>
                        <Label for='jobType'>Select the type of message to send</Label>
                        <Input id='jobType' value={this.state.type} onChange={this.handleChangeType} type="select">
                            {Object.keys(Type).map( i => {
                                return <option key={i} value={i}>{i.charAt(0)+i.substring(1).toLowerCase()}</option>
                            })}
                        </Input>
                    </FormGroup>
                    <FormGroup>
                        <Label for='jobMessage'>Message to send</Label>
                        <Input onChange={this.handleChangeMessage} type="text" id="jobMessage" value={this.state.message} required/>
                    </FormGroup>
                    {this.state.type === "JOB" && 
                    <div>
                    <FormGroup>
                        <Label for='material'>Select the material to use</Label>
                        <Input id='material' value={this.state.material} onChange={this.handleChangeMaterial} type="select">
                            {Object.keys(Material).map( i => {
                                return <option key={i} value={i}>{i}</option>
                            })}
                        </Input>
                    </FormGroup>
                    <FormGroup>
                        <Label for="IPAddress">Device IP Address</Label>
                        <Input id="IPAddress" value={this.state.IPAddress} onChange={this.handleChangeIP}
                                type="text" valid={this.validateIP()} invalid={!this.validateIP()}/>
                    </FormGroup>
                    <FormGroup>
                        <Label for="file">File</Label>
                        <Input id="file" name="file" type="file" 
                            accept='.pdf, .stl' 
                            onChange={this.handleChangeFile} required/>
                        <FormText>Select the file to send to the job</FormText>
                    </FormGroup>
                    </div>
                    }

                    <FormGroup>
                        <Input type="submit" value="Submit" />
                    </FormGroup>
                </FormRS>

                {/* Renders a modal stating that the new job is being processed whenever a new now has been submitted until
                    response from server is received */}
                {this.state.processingJob &&
                <Modal centered isOpen={true}>
                    <ModalHeader>Processing</ModalHeader>
                    <ModalBody> 
                        <Spinner/>
                        {' '}
                        Your job is being sent, please wait
                    </ModalBody>
                </Modal>
                }

                {/* Renders a modal to inform about last job submission outcome*/}
                {this.state.submitOutcome !== '' &&
                <Modal centered isOpen={true}>
                    <ModalHeader>Outcome</ModalHeader>
                    <ModalBody> 
                        {this.state.submitOutcome}
                    </ModalBody>
                    <ModalFooter>
                        <Button
                            color="primary"
                            onClick={this.resetForm}>
                            OK!
                        </Button>
                    </ModalFooter>
                </Modal>

                }
            </div>
        )
    }
}

export default Form;