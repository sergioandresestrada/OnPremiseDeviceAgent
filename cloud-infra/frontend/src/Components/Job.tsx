import React, { FormEvent } from "react";
import { Link } from "react-router-dom";
import { Form as FormRS, FormGroup, Input, Label, Modal, ModalBody, ModalHeader, Spinner, Button, ModalFooter, FormText, Alert} from 'reactstrap';
import { DevicePublic } from "../utils/types";
import { isValidFile, URL } from '../utils/utils';
import Help from "./Help";

enum Material {
    "HR PA 11" = "HR PA 11", 
    "HR PA 12GB" = "HR PA 12GB",
    "HR TPA" = "HR TPA", 
    "HR PP" = "HR PP", 
    "HR PA 12" = "HR PA 12"
}

interface IJob{
    file? : File,
    material? : Material,
    availableDevices : DevicePublic[],
    selectedDeviceName : string

    processingJob : boolean,
    submitOutcome : string
}

const initialState = {
    file : undefined,
    material : Material['HR PA 11'],
    availableDevices : [] as DevicePublic[],
    selectedDeviceName : '',


    processingJob : false,
    submitOutcome : ''
}

class Job extends React.Component<{}, IJob>{
    constructor(props: any){
        super(props)
        this.state = initialState

        this.handleSubmit = this.handleSubmit.bind(this)
        this.handleChangeFile = this.handleChangeFile.bind(this)
        this.handleChangeMaterial = this.handleChangeMaterial.bind(this)
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
        
        let fullURL : string = ""
        let fetchOptions: object = {}

        if (this.state.file === undefined){
            alert("Error getting the file")
            return
        }
        if(!isValidFile(this.state.file)){
            alert("Invalid file selected. Please select a PDF or STL file and try again")
            return
        }

        fullURL = URL + "/job"
        let formData = new FormData()
        let data = JSON.stringify({
            type: "JOB",
            material: this.state.material,
            DeviceName: this.state.selectedDeviceName
        })
        formData.append("data", data)
        formData.append("file", this.state.file)

        fetchOptions = {
            method: "POST",
            body: formData
        }

        this.setState({
            processingJob : true
        })

        fetch(fullURL, fetchOptions)
        .then(response => {
            let outcome = ""
            switch(response.status){
                case 200:
                    outcome = "New Job was sent successfully."
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

    handleChangeFile(event : React.ChangeEvent<HTMLInputElement>){
        if (event.target.files != null){
            this.setState({
                file : event.target.files[0]
            })
        }
    }

    handleChangeMaterial(event: React.ChangeEvent<HTMLInputElement>){
        this.setState({
            material : event.target.value as Material
        })
    }

    handleChangeSelectedDevice(event : React.ChangeEvent<HTMLInputElement>){
        this.setState({
            selectedDeviceName : event.target.value
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
                        <Label for='material'>Select the material to use</Label>
                        <Input id='material' value={this.state.material} onChange={this.handleChangeMaterial} type="select">
                            {Object.keys(Material).map( i => {
                                return <option key={i} value={i}>{i}</option>
                            })}
                        </Input>
                    </FormGroup>
                    <FormGroup>
                        <Label for="file">File</Label>
                        <Input id="file" name="file" type="file" 
                            accept='.pdf, .stl' 
                            onChange={this.handleChangeFile} required/>
                        <FormText>Select the file to send to the job</FormText>
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
                                <Button type="submit" color="primary" outline style={{width:"100%"}}>Print</Button>
                            </FormGroup>
                        </div>
                    }
                </FormRS>

                <Help 
                    message={"You can use a Job message to print an archive in the desired device.\n"+
                                "Select the desired material to be used, the device and the corresponding file (only STL and PDF are supported)\n"+
                                "and click the 'Print' button."} 
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

export default Job;