import React from "react";
import { Link } from "react-router-dom";
import { Alert, Button, Modal, ModalBody, ModalFooter, ModalHeader, Spinner } from "reactstrap";
import { URL, beautifyFileName } from "../utils/utils"
import '../App.css';
import ShowIdentification from "./DeviceIdentificationRender";
import ShowJobs from "./DeviceJobsRender";

interface IRequestedDeviceInfo{
    errorInFetch: boolean,
    isLoading: boolean,
    information: any,
    errorMessage: string
}

const initialState = {
    errorInFetch: false,
    isLoading: true,
    information: {},
    errorMessage: ""
}

class RequestedDeviceInfo extends React.Component<{},IRequestedDeviceInfo>{
    constructor(props: any){
        super(props)
        this.state = initialState
    }


    componentDidMount(){
        if (localStorage.getItem("requestedInfo") === null ){
            this.setState({
                isLoading: false,
                errorInFetch: true,
                errorMessage: "You have not requested any information. " +
                    "Please choose the desired information from the following list"
            })
            return
        }

        let fullURL : string = ""

        fullURL = URL + "/getInformationFile?file=" + localStorage.getItem("requestedInfo")
        //fullURL = URL+"/testjobs"
        fetch(fullURL)
        .then(res => res.json())
        .then(
            (result) => {
                this.setState({
                    information: result as any,
                    isLoading: false,
                    errorInFetch: false
                })
            },
            (error) => {
                this.setState({
                    errorInFetch: true,
                    errorMessage: "The requested file is no longer available or there was a server error. Please, try again later",
                    isLoading: false
                })
            }
        )
    }

    render(){
        const {errorInFetch, isLoading } = this.state
        const fileName = localStorage.getItem("requestedInfo")
        if (isLoading){
            return(
                <Modal centered isOpen={true}>
                    <ModalHeader>Getting data</ModalHeader>
                    <ModalBody> 
                        <Spinner/>
                        {' '}
                        Requested information is getting loaded
                    </ModalBody>
                </Modal>
            )
        }
        if (errorInFetch) {
            return (
                <Modal centered isOpen={true}>
                <ModalHeader>Error</ModalHeader>
                <ModalBody> 
                    {this.state.errorMessage}
                </ModalBody>
                <ModalFooter>
                    <Link to="/deviceInfoList" style={{ color: "#0096D6", textDecoration: "none" }}>OK!</Link>
                </ModalFooter>
            </Modal>
            )
        }      

        if(this.state.information.hasOwnProperty("Jobs")){
            return(
                <div className="DeviceInfoShow">
                    <Alert color="primary">
                        {fileName !== null ? "Loaded information: Jobs from device " + beautifyFileName(fileName) : ""}
                    </Alert>
                    <ShowJobs Jobs={this.state.information.Jobs} />
                    <Button style={{ marginTop: '2rem', backgroundColor:"#0096D6"}}>
                        <Link to="/deviceInfoList" style={{color:"white", textDecoration: "none" }}>Go back to available information List</Link>
                    </Button>
                </div>
            )
        }
        
        if(this.state.information.hasOwnProperty("Identification")){
            
            return(
                <div className="DeviceInfoShow">
                    <Alert color="primary">
                        {fileName !== null ? "Loaded information: Identification from device " + beautifyFileName(fileName) : ""}
                    </Alert>
                    <ShowIdentification Identification={this.state.information.Identification}/>
                    <Button style={{ marginTop: '2rem', backgroundColor:"#0096D6"}}>
                        <Link to="/deviceInfoList" style={{color:"white", textDecoration: "none" }}>Go back to available information List</Link>
                    </Button>
                </div>
            )
        }
    }
}

export default RequestedDeviceInfo