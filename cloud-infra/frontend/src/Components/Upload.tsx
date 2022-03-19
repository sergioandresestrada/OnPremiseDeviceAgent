import React, { FormEvent } from "react";
import { Form as FormRS, FormGroup, Input, Label, Modal, ModalBody, ModalHeader, Spinner, Button, ModalFooter} from 'reactstrap';
import { URL, validateIP } from '../utils/utils';

enum UploadInfoTypes {
    "Jobs" = "Jobs",
    "Identification" = "Identification"
}

interface IUpload {
    IPAddress? : string,
    UploadInfo? : UploadInfoTypes,

    processingJob : boolean,
    submitOutcome : string
}

const initialState = {
    IPAddress : "",
    UploadInfo : UploadInfoTypes["Jobs"],

    processingJob : false,
    submitOutcome : ''
}

class Upload extends React.Component<{}, IUpload>{
    constructor(props: any){
        super(props)
        this.state = initialState

        this.handleSubmit = this.handleSubmit.bind(this)
        this.handleChangeIP = this.handleChangeIP.bind(this)
        this.handleChangeUploadInfo = this.handleChangeUploadInfo.bind(this)
    }

    handleSubmit(event : FormEvent){
        event.preventDefault()

        let fetchOptions : object = {}
        let fullURL : string = ""

        if(!validateIP(this.state.IPAddress)){
            alert("Invalid IP address. Check it and try again.")
            return
        }
        fullURL = URL + "/upload"
        fetchOptions = {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                type: "UPLOAD",
                IPAddress: this.state.IPAddress,
                UploadInfo : this.state.UploadInfo
            })
        }

        this.setState({
            processingJob : true
        })

        fetch(fullURL, fetchOptions)
        .then(response => {
            let outcome = ""
            switch(response.status){
                case 200:
                    outcome = "New Upload request was sent successfully."
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
    
    handleChangeIP(event: React.ChangeEvent<HTMLInputElement>){
        this.setState({
            IPAddress : event.target.value
        })
    }

    handleChangeUploadInfo(event: React.ChangeEvent<HTMLInputElement>){
        this.setState({
            UploadInfo : event.target.value as UploadInfoTypes
        })
    }

    resetForm = () => {
        this.setState(initialState)
    }

    render() {
        return(
            <div>
                <FormRS onSubmit={this.handleSubmit}>
                    <FormGroup>
                        <Label for="IPAddress">Device IP Address</Label>
                        <Input id="IPAddress" value={this.state.IPAddress} onChange={this.handleChangeIP}
                                type="text" valid={validateIP(this.state.IPAddress)} invalid={!validateIP(this.state.IPAddress)}/>
                    </FormGroup>

                    <FormGroup>
                        <Label for='uplaodInfo'>Select the information to request</Label>
                        <Input id='uplaodInfo' value={this.state.UploadInfo} onChange={this.handleChangeUploadInfo} type="select">
                            {Object.keys(UploadInfoTypes).map( i => {
                                return <option key={i} value={i}>{i}</option>
                            })}
                        </Input>
                    </FormGroup>
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

export default Upload;