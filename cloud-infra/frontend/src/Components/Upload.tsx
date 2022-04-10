import React, { FormEvent } from "react";
import { Link } from "react-router-dom";
import { Form as FormRS, FormGroup, Input, Label, Modal, ModalBody, ModalHeader, Spinner, Button, ModalFooter, Alert} from 'reactstrap';
import { DevicePublic } from "../utils/types";
import { URL } from '../utils/utils';
import Help from "./Help";

enum UploadInfoTypes {
    "Jobs" = "Jobs",
    "Identification" = "Identification"
}

interface IUpload {
    availableDevices : DevicePublic[],
    selectedDeviceName : string,
    UploadInfo? : UploadInfoTypes,

    processingJob : boolean,
    submitOutcome : string
}

const initialState = {
    availableDevices : [] as DevicePublic[],
    selectedDeviceName : '',
    UploadInfo : UploadInfoTypes["Jobs"],

    processingJob : false,
    submitOutcome : ''
}

class Upload extends React.Component<{}, IUpload>{
    constructor(props: any){
        super(props)
        this.state = initialState

        this.handleSubmit = this.handleSubmit.bind(this)
        this.handleChangeUploadInfo = this.handleChangeUploadInfo.bind(this)
        this.handleChangeSelectedDevice = this.handleChangeSelectedDevice.bind(this)
    }

    componentDidMount(){
        fetch(URL + "/getPublicDevices")
        .then(res => res.json())
        .then(
            (result) => {
                this.setState({
                    availableDevices: result as DevicePublic[]
                })
                if (this.state.availableDevices.length > 0){
                    this.setState({
                        selectedDeviceName : this.state.availableDevices[0].Name
                    })
                }
            }
        )
        .catch(error => {
            alert("There was an error connecting to the server, please try again later.")
        })
    }

    handleSubmit(event : FormEvent){
        event.preventDefault()

        let fetchOptions : object = {}
        let fullURL : string = ""

        fullURL = URL + "/upload"
        fetchOptions = {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                type: "UPLOAD",
                DeviceName: this.state.selectedDeviceName,
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
    
    handleChangeSelectedDevice(event : React.ChangeEvent<HTMLInputElement>){
        this.setState({
            selectedDeviceName : event.target.value
        })
    }

    handleChangeUploadInfo(event: React.ChangeEvent<HTMLInputElement>){
        this.setState({
            UploadInfo : event.target.value as UploadInfoTypes
        })
    }

    resetForm = () => {
        let copyDevices = this.state.availableDevices
        this.setState(initialState)
        this.setState({
            availableDevices : copyDevices
        })
    }

    render() {
        return(
            <div>
                <FormRS onSubmit={this.handleSubmit}>
                    <FormGroup>
                        <Label for='uplaodInfo'>Select the information to request</Label>
                        <Input id='uplaodInfo' value={this.state.UploadInfo} onChange={this.handleChangeUploadInfo} type="select">
                            {Object.keys(UploadInfoTypes).map( i => {
                                return <option key={i} value={i}>{i}</option>
                            })}
                        </Input>
                    </FormGroup>

                    {this.state.availableDevices.length === 0 &&
                        <Alert color="warning">No devices available! <Link to="/devices/new" style={{ color: "#0096D6", textDecoration: "none"}}>Go add some now.</Link></Alert>
                    }
                    {this.state.availableDevices.length !== 0 &&
                        <div>
                            <FormGroup>
                                <Label for='device'>Select the device</Label>
                                <Input id='device' value={this.state.selectedDeviceName} onChange={this.handleChangeSelectedDevice} type="select">
                                    {this.state.availableDevices.map((dev)  => {
                                        return <option key={dev.Name} value={dev.Name}>{dev.Name + (dev.Model === undefined ? "" : (" - " + dev.Model))}</option>
                                    })}
                                </Input>            
                            </FormGroup>
                            <FormGroup>
                                <Button type="submit" color="primary" outline style={{width:"100%"}}> {"Request " + this.state.UploadInfo + " information"}</Button>
                            </FormGroup>
                        </div>
                    }
                </FormRS>
                
                <Help 
                    message={"You can use an Upload message to request information from a device.\n"+
                                "Requested information can be device's Jobs list or device's identification.\n"+
                                "Select the device an desired information to request and click the button.\n"+
                                "Requested information can be checked, when received back, on the corresponding tab you can find at the top"} 
                    opened={false}
                />

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