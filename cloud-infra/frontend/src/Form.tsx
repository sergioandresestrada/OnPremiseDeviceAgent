
import React, { FormEvent } from 'react';
import { Form as FormRS, FormGroup, FormText, Input, Label, Modal, ModalBody, ModalHeader, Spinner, Button, ModalFooter} from 'reactstrap';

enum Type {
    HEARTBEAT = "HEARTBEAT",
    JOB = "JOB",
    PLACEHOLDER1 = "PLACEHOLDER1"
}

const initialState = {
    message : '',
    type : Type.HEARTBEAT,
    file : undefined,

    processingJob : false,
    submitOutcome : ''
}

const URL = "http://192.168.1.208:12345"

interface IJob{
    message : string,
    type : Type,
    file? : File

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
                    message: this.state.message
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
            if (response.status === 200){
                outcome = "New job was sent successfully"
            } else {
                outcome = "There was a problem sending the new job"
            }
            this.setState({
                submitOutcome : outcome,
                processingJob : false
            })
        })
        .catch(error => {
            alert(error)
            let outcome = "There was an error processing the petition, please check the fields and try again"
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
        if (!(acceptedTypes.indexOf(fileExtension) > -1)){
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
                    <FormGroup>
                        <Label for="file">File</Label>
                        <Input id="file" name="file" type="file" 
                            accept='.pdf, .stl' 
                            onChange={this.handleChangeFile} required/>
                        <FormText>Select the file to send to the job</FormText>
                    </FormGroup>
                    }
                    {/* application/pdf,application/x-pdf,application/x-bzpdf,application/x-gzpdf,*/ }
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