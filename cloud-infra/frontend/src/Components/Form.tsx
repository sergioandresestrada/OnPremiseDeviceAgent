
import React from 'react';
import Heartbeat from './Heartbeat';
import Job from './Job';
import Upload from './Upload';
import { Alert, Col, Container, Form as FormRS, FormGroup, Input, Label, Modal, ModalBody, ModalHeader, Row, Spinner} from 'reactstrap';
import '../App.css';
import { DevicePublic } from '../utils/types';
import { URL } from '../utils/utils';
import { Link } from 'react-router-dom';

enum Type {
    HEARTBEAT = "HEARTBEAT",
    JOB = "JOB",
    UPLOAD = "UPLOAD"
}

const initialState = {
    type : Type.HEARTBEAT,
    availableDevices : [] as DevicePublic[],
    errorInFetch: false,
    isLoading: true,
}


interface IJob{
    type : Type,
    availableDevices : DevicePublic[],
    errorInFetch: boolean,
    isLoading: boolean,
}

class Form extends React.Component<{}, IJob>{ 
    constructor(props : any){
        super(props);
        this.state = initialState
        this.handleChangeType = this.handleChangeType.bind(this)
    }

    componentDidMount(){
        fetch(URL + "/getPublicDevices")
        .then(res => res.json())
        .then(
            (result) => {
                this.setState({
                    availableDevices: result as DevicePublic[],
                    isLoading : false,
                    errorInFetch : false
                })
            }
        )
        .catch(error => {
            this.setState({
                isLoading : false,
                errorInFetch : true
            })
        })
    }

    
    handleChangeType(event : React.ChangeEvent<HTMLInputElement>){
        this.setState({
            type : event.target.value as Type
        })
    }

    render() {
        const {errorInFetch, isLoading } = this.state

        if (isLoading){
            return (
                <div>
                    <Modal centered isOpen={true}>
                        <ModalHeader>Getting data</ModalHeader>
                        <ModalBody> 
                            <Spinner/>
                            {' '}
                            Available devices are getting loaded
                        </ModalBody>
                    </Modal>
                </div>
            )
        }

        if (errorInFetch){
            return(
                <Modal centered isOpen={true}>
                    <ModalHeader>Error :(</ModalHeader>
                    <ModalBody> 
                        Server is down. Please try again later.
                    </ModalBody>
                </Modal>
            )
        }

        return(
            <Container>
                <Row>
                    <Col className='Form' style={{maxWidth:"600px"}}>
                    <FormRS>
                        <FormGroup>
                            <Label for='jobType'>Select the type of message to send</Label>
                            <Input id='jobType' value={this.state.type} onChange={this.handleChangeType} type="select">
                                {Object.keys(Type).map( i => {
                                    return <option key={i} value={i}>{i.charAt(0)+i.substring(1).toLowerCase()}</option>
                                })}
                            </Input>
                        </FormGroup>
                    </FormRS>
                    {this.state.availableDevices.length === 0 &&
                        <Alert color="warning">No devices available! <Link to="/devices/new" style={{ color: "#0096D6", textDecoration: "none"}}>Go add some now.</Link></Alert>
                    }
                    {this.state.availableDevices.length > 0 &&
                        <div>
                            {this.state.type === "HEARTBEAT" &&
                                <Heartbeat devices={this.state.availableDevices}/>
                            }

                            {this.state.type === "JOB" && 
                                <Job devices={this.state.availableDevices}/>
                            }

                            {this.state.type === "UPLOAD" &&
                                <Upload devices={this.state.availableDevices}/>
                            }
                        </div>
                    }
                    </Col>
                </Row>
            </Container>
        )
    }
}

export default Form;