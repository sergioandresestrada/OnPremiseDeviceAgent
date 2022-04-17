
import React, { FormEvent } from "react";
import { Link } from "react-router-dom";
import { Button, Form as FormRS, FormGroup, Input, Label, Modal, ModalBody, ModalFooter, ModalHeader, Spinner } from 'reactstrap';
import { Device } from "../utils/types";
import { URL, validateIP } from "../utils/utils";

interface PDeviceForm {
    isNewDevice: boolean,
}

interface IDeviceForm{
    processing: boolean,
    submitOutcome : string,

    editDeviceUUID?: string,

    isLoading: boolean,
    errorInFetch: boolean,

    deviceName: string,
    deviceIP: string,
    deviceModel: string
}

const initialState = {
    processing : false,
    submitOutcome : '',

    editDeviceUUID: undefined,

    isLoading: true,
    errorInFetch: false,

    deviceName : '',
    deviceIP: '',
    deviceModel : ''
}

class DeviceForm extends React.Component<PDeviceForm,IDeviceForm> {
    constructor(props: any){
        super(props)
        this.state = initialState

        this.handleChangeIP = this.handleChangeIP.bind(this)
        this.handleChangeName = this.handleChangeName.bind(this)
        this.handleChangeModel = this.handleChangeModel.bind(this)
        this.handleSubmitNew = this.handleSubmitNew.bind(this)
        this.handleSubmitUpdate = this.handleSubmitUpdate.bind(this)
    }  

    componentDidMount(){
        if (!this.props.isNewDevice){
            let path = window.location.pathname.split("/")
            let deviceUUID = path[path.length-1]
            this.setState({
                editDeviceUUID : deviceUUID
            })

            fetch(URL + "/devices/" + deviceUUID)
            .then(res => res.json())
            .then(
                (result) => {
                    let device = result as Device
                    this.setState({
                        deviceName : device.Name,
                        deviceModel : device.Model,
                        deviceIP: device.IP,
                        isLoading: false,
                        errorInFetch: false                   
                    })
                }
            )
            .catch(error => {
                this.setState({
                    isLoading : false,
                    errorInFetch : true,
                })
            })
        }
        else {
            this.setState({
                isLoading: false,
            })
        }
    }

    handleChangeName(event: React.ChangeEvent<HTMLInputElement>){
        this.setState({
            deviceName : event.target.value
        })
    }

    handleChangeModel(event: React.ChangeEvent<HTMLInputElement>){
        this.setState({
            deviceModel : event.target.value
        })
    }

    handleChangeIP(event: React.ChangeEvent<HTMLInputElement>){
        this.setState({
            deviceIP : event.target.value
        })
    }

    handleSubmitUpdate(event: FormEvent){
        event.preventDefault()

        event.preventDefault()
        
        let fullURL = URL + "/devices/" + this.state.editDeviceUUID
        let fetchOptions : object = {}
        let body : object = {}

        if (this.state.deviceModel.trim() !== ""){
            body = {
                DeviceUUID : this.state.editDeviceUUID,
                Model : this.state.deviceModel,
                Name : this.state.deviceName,
                IP: this.state.deviceIP
            }
        }
        else {
            body = {
                DeviceUUID : this.state.editDeviceUUID,
                Name : this.state.deviceName,
                IP: this.state.deviceIP
            }
        }
        fetchOptions = {
            method: "PUT",
            headers : {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(body)
        }

        this.setState({
            processing : false
        })

        fetch(fullURL, fetchOptions)
        .then(response => {
            let outcome = ""
            switch(response.status){
                case 200:
                    outcome = "Device was updated successfully"
                    break
                case 400:
                    outcome = "Bad request. Check the fields and try again. Remember that device name and IP has to be unique."
                    break
                case 500:
                    outcome = "Server error. Try again later."
                    break
            }
            this.setState({
                submitOutcome : outcome,
                processing : false
            })
        })
        .catch(error => {
            let outcome = "There was an error connecting to the server, please try again later."
            this.setState({
                submitOutcome : outcome,
                processing : false
            })
        })

    }

    handleSubmitNew(event: FormEvent){
        event.preventDefault()
        
        let fullURL = URL + "/devices"
        let fetchOptions : object = {}
        let body : object = {}

        if (this.state.deviceModel.trim() !== ""){
            body = {
                Model : this.state.deviceModel,
                Name : this.state.deviceName,
                IP: this.state.deviceIP
            }
        }
        else {
            body = {
                Name : this.state.deviceName,
                IP: this.state.deviceIP
            }
        }
        fetchOptions = {
            method: "POST",
            headers : {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(body)
        }

        this.setState({
            processing : false
        })

        fetch(fullURL, fetchOptions)
        .then(response => {
            let outcome = ""
            switch(response.status){
                case 200:
                    outcome = "New Device was added successfully"
                    break
                case 400:
                    outcome = "Bad request. Check the fields and try again. Remember that device name and IP has to be unique."
                    break
                case 500:
                    outcome = "Server error. Try again later."
                    break
            }
            this.setState({
                submitOutcome : outcome,
                processing : false
            })
        })
        .catch(error => {
            let outcome = "There was an error connecting to the server, please try again later."
            this.setState({
                submitOutcome : outcome,
                processing : false
            })
        })
    }

    render(){
        const { errorInFetch, isLoading } = this.state

        /* Renders a modal while the device is being requested */
        if (isLoading){
            return (
                <div>
                    <Modal centered isOpen={true}>
                        <ModalHeader>Getting data</ModalHeader>
                        <ModalBody> 
                            <Spinner/>
                            {' '}
                            Your device is getting loaded
                        </ModalBody>
                    </Modal>
                </div>
            )
        }
        
        /* Renders a modal indicating that there was an error while requesting the device */
        if (errorInFetch){
            return(
            <Modal centered isOpen={true}>
                <ModalHeader>Error</ModalHeader>
                <ModalBody> 
                    There was an error while requesting your device. Please try again later
                </ModalBody>
                <ModalFooter>
                    <Link to="/" style={{ color: "#0096D6", textDecoration: "none" }}>OK!</Link>
                </ModalFooter>
            </Modal>
            )
        }

        return (
            <div className='Form'>
                <h5 style={{textAlign: 'center', color:'#0096D6', paddingBottom:'1em'}}>{(this.props.isNewDevice ? 'New' : 'Update') + " Device" }</h5>
                
                <FormRS onSubmit={this.props.isNewDevice? this.handleSubmitNew : this.handleSubmitUpdate}>
                    <FormGroup>
                        <Label for='deviceName'>Device Name</Label>
                        <Input onChange={this.handleChangeName} type="text" id="deviceName" value={this.state.deviceName} required/>
                    </FormGroup>
                    <FormGroup>
                        <Label for="deviceIP">Device IP Address</Label>
                        <Input id="deviceIP" value={this.state.deviceIP} onChange={this.handleChangeIP}
                                type="text" valid={validateIP(this.state.deviceIP)} invalid={!validateIP(this.state.deviceIP)}/>
                    </FormGroup>
                    <FormGroup>
                        <Label for='deviceModel'>Device Model</Label>
                        <Input onChange={this.handleChangeModel} type="text" id="deviceModel" value={this.state.deviceModel}/>
                    </FormGroup>
                    <FormGroup>
                        <Button type="submit" color="primary" outline style={{width:"100%"}}> {"Save " + (this.props.isNewDevice ? "new device" : "changes")}</Button>
                    </FormGroup>
                </FormRS>


                {/* Renders a modal stating that the device is being processed until
                    response from server is received */}
                {this.state.processing &&
                <Modal centered isOpen={true}>
                    <ModalHeader>Processing</ModalHeader>
                    <ModalBody> 
                        <Spinner/>
                        {' '}
                        Your Device is being processed, please wait
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
                        <Link to="/devices" style={{ color: "#0096D6", textDecoration: "none" }}>OK!</Link>
                    </ModalFooter>
                </Modal>
                }
            </div>
        )
    }
    
}

export default DeviceForm