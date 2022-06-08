import React from "react";
import { Link, Navigate } from "react-router-dom";
import { Alert, Button, Container, Modal, ModalBody, ModalFooter, ModalHeader, Spinner, Table } from "reactstrap";
import { Device } from "../utils/types";
import { URL } from "../utils/utils"

interface IDeviceList {
    devices: Device[]
    errorInFetch: boolean,
    errorBody: string
    isLoading: boolean,

    redirectEdit: boolean,
    redirectMessages: boolean,
    DeviceUUID: string
}

const initialState = {
    devices: [],
    errorInFetch: false,
    errorBody: "",
    isLoading: true,

    redirectEdit: false,
    redirectMessages: false,
    DeviceUUID: ''
}

class DeviceList extends React.Component<{},IDeviceList> {
    constructor(props: any){
        super(props)
        this.state = initialState
        
        this.loadDevices = this.loadDevices.bind(this)
        this.deleteDevice = this.deleteDevice.bind(this)
        this.redirectToEdit = this.redirectToEdit.bind(this)

    }

    componentDidMount(){
        this.setState(initialState)
        this.loadDevices()
    }

    loadDevices(){
        fetch(URL + "/devices")
        .then(res => res.json())
        .then(
            (result) => {
                this.setState({
                    devices: result as Device[],
                    isLoading: false,
                    errorInFetch: false                   
                })
            }
        )
        .catch(error => {
            this.setState({
                isLoading : false,
                errorInFetch : true,
                errorBody: "There was an error while requesting the available information. Please try again later."
            })
        })
    }

    deleteDevice(uuid: string) {
        fetch(URL + "/devices/"+uuid, {
            method: "DELETE"
        })
        .then(res => {
            this.loadDevices()
        }
        )
        .catch(error => {
            this.setState({
                errorInFetch: true,
                errorBody: "Error while deleting the selected device, try again later"
            })
        })
    }

    redirectToEdit(uuid: string){
        this.setState({
            redirectEdit: true,
            DeviceUUID: uuid
        })
    }

    redirectToMessages(uuid: string){
        this.setState({
            redirectMessages: true,
            DeviceUUID: uuid
        })
    }


    render(){
        const { errorInFetch, isLoading, redirectEdit, redirectMessages } = this.state
        
        if (redirectEdit){
            return (
                <Navigate to={"/devices/edit/"+this.state.DeviceUUID}></Navigate>
            )
        }
        
        if (redirectMessages){
            return (
                <Navigate to={"/messages/"+this.state.DeviceUUID}></Navigate>
            )
        }

        /* Renders a modal while the devices are being requested */
        if (isLoading){
            return (
                <div>
                    <Modal centered isOpen={true}>
                        <ModalHeader>Getting data</ModalHeader>
                        <ModalBody> 
                            <Spinner/>
                            {' '}
                            The list of devices is getting loaded
                        </ModalBody>
                    </Modal>
                </div>
            )
        }
        
        /* Renders a modal indicating that there was an error while requesting the informatino */
        if (errorInFetch){
            return(
            <Modal centered isOpen={true}>
                <ModalHeader>Error</ModalHeader>
                <ModalBody> 
                    {this.state.errorBody}
                </ModalBody>
                <ModalFooter>
                    <Link to="/" style={{ color: "#0096D6", textDecoration: "none" }}>OK!</Link>
                </ModalFooter>
            </Modal>
            )
        }

        // When there are no devices, we show an alert and a link to the insert form
        if (this.state.devices.length === 0){
            return(
                <Container style={{marginTop: '2rem'}}>
                    <Alert color="warning">No devices available! <Link to="/devices/new" style={{ color: "#0096D6", textDecoration: "none"}}>Go add some now.</Link></Alert>
                </Container>
            )
        }

        return(
            <div>
                <Container style={{marginTop: '2rem'}} className="custom-background">
                <Button style={{marginBottom:'1.5em', backgroundColor:"#0096D6", width: "175px"}} tag={Link} to="/devices/new">Add new Device</Button>
                <Table hover>
                    <thead>
                        <tr>
                            <th>Name</th>
                            <th>IP Address</th>
                            <th>Model</th>
                            <th>Last message result</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                    {Object.values(this.state.devices).map((entry) => {
                        return(
                            <tr key={entry.DeviceUUID}>
                                <td>{entry.Name}</td>
                                <td>{entry.IP}</td>
                                <td>{entry.Model}</td>
                                <td>
                                    {entry.LastResult === undefined ? "Unknown" : entry.LastResult.substring(0,7) }
                                </td>
                                <td>
                                    <Button color="primary" onClick={() => this.redirectToEdit(entry.DeviceUUID)} outline>Edit</Button>
                                    <Button color="danger" onClick={() => this.deleteDevice(entry.DeviceUUID)} outline style={{marginLeft: "1em"}}>Delete</Button>
                                    <Button color="success" onClick={() => this.redirectToMessages(entry.DeviceUUID)} outline style={{marginLeft: "1em"}}>Messages</Button>
                                </td>
                            </tr>
                        )
                    })}
                    </tbody>
                </Table>
                </Container>
            </div>
        )
    }
}

export default DeviceList;