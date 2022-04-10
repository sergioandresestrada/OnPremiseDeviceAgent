import React, { FormEvent } from "react";
import { Link } from "react-router-dom";
import { Form as FormRS, FormGroup, Input, Label, Modal, ModalBody, ModalHeader, Spinner, Button, ModalFooter, Alert} from 'reactstrap';
import { DevicePublic } from "../utils/types";
import { URL } from '../utils/utils';
import Help from "./Help";

interface IHeartbeat{
    message : string,
    availableDevices : DevicePublic[],
    selectedDeviceName : string

    processingHB : boolean,
    submitOutcome : string
}

const initialState = {
    message : '',
    availableDevices : [] as DevicePublic[],
    selectedDeviceName : '',

    processingHB : false,
    submitOutcome : ''
}

class Heartbeat extends React.Component<{}, IHeartbeat>{
    constructor(props: any){
        super(props)
        this.state = initialState
        
        this.handleSubmit = this.handleSubmit.bind(this)
        this.handleChangeMessage = this.handleChangeMessage.bind(this)
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

    handleChangeMessage(event: React.ChangeEvent<HTMLInputElement>){
        this.setState({
            message : event.target.value
        });
    }

    handleChangeSelectedDevice(event : React.ChangeEvent<HTMLInputElement>){
        this.setState({
            selectedDeviceName : event.target.value
        })
    }

    handleSubmit(event: FormEvent){
        event.preventDefault()
        
        let fullURL : string = ""
        let fetchOptions: object = {}

        fullURL = URL + "/heartbeat"
        fetchOptions = {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                type: "HEARTBEAT",
                message: this.state.message,
                DeviceName : this.state.selectedDeviceName
            })
        }

        this.setState({
            processingHB : true
        })

        fetch(fullURL, fetchOptions)
        .then(response => {
            let outcome = ""
            switch(response.status){
                case 200:
                    outcome = "New Heartbeat was sent successfully."
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
                processingHB : false
            })
        })
        .catch(error => {
            let outcome = "There was an error connecting to the server, please try again later."
            this.setState({
                submitOutcome : outcome,
                processingHB : false
            })
        })
    }

    resetForm = () => {
        let copyDevices = this.state.availableDevices
        this.setState(initialState)
        this.setState({
            availableDevices : copyDevices
        })
    }

    render(){
        return(
            <div>
                <FormRS onSubmit={this.handleSubmit}>
                    <FormGroup>
                        <Label for='heartbeatMessage'>Message to send</Label>
                        <Input onChange={this.handleChangeMessage} type="text" id="heartbeatMessage" value={this.state.message} required/>
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
                                <Button type="submit" color="primary" outline style={{width:"100%"}}> Send Heartbeat to the printer</Button>
                            </FormGroup>
                        </div>
                    }
                </FormRS>

                <Help 
                    message={"You can use a Heartbeat message to test whether a device is active or not.\n"+
                                "Introduce the desired message to send, select a device and click the button."} 
                    opened={false}
                />
              
                {/* Renders a modal stating that the new HB is being processed whenever a new one has been submitted until
                    response from server is received */}
                {this.state.processingHB &&
                <Modal centered isOpen={true}>
                    <ModalHeader>Processing</ModalHeader>
                    <ModalBody> 
                        <Spinner/>
                        {' '}
                        Your Heartbeat is being sent, please wait
                    </ModalBody>
                </Modal>
                }

                {/* Renders a modal to inform about last HB submission outcome*/}
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

export default Heartbeat;
