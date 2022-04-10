import React, { ReactNode } from "react";
import { Link, Navigate } from "react-router-dom";
import {Alert, Button, Col, Container, Modal, ModalBody, ModalFooter, ModalHeader, Row, Spinner, Table} from "reactstrap";
import { URL } from "../utils/utils"

interface IDeviceInfoList{
    jobs: string[],
    identification: string[],
    errorInFetch: boolean,
    isLoading: boolean,
    redirect: boolean
}

const initialState = {
    jobs: [],
    identification: [],
    errorInFetch: false,
    isLoading: true,
    redirect: false
}

class DeviceInfoList extends React.Component<{},IDeviceInfoList>{

    constructor(props: any){
        super(props)
        this.state = initialState
        
        this.renderRows = this.renderRows.bind(this)
        this.selectDeviceInfo = this.selectDeviceInfo.bind(this)
    }

    componentDidMount(){
        let fullURL : string = ""

        fullURL = URL + "/availableInformation"
        fetch(fullURL)
        .then(res => res.json())
        .then(
            (result) => {
                this.setState({
                    jobs: result.Jobs,
                    identification: result.Identification,
                    isLoading: false,
                    errorInFetch: false
                })
            },
            (error) => {
                this.setState({
                    errorInFetch: true,
                    isLoading: false
                })
            }
        )
    }

    renderRows(data: string[]): ReactNode{
        if (data.length === 0) {
            return (
                <tr key="empty" style={{color:"black"}}>
                    <td>No available information at the moment</td>
                </tr>
            )
        }
        return data.map((item) => {
            return(
                <tr key={item} id={item} onClick={this.selectDeviceInfo} style={{color:"black"}}>
                    <td>{item}</td>
                </tr>
            )
        })
    }

    selectDeviceInfo(e: React.MouseEvent<HTMLTableRowElement>){
        const selectedRow = e.target as HTMLTableRowElement        
        localStorage.setItem("requestedInfo", selectedRow.innerText)
        this.setState({
            redirect: true
        })
    }

    render(){
        const { errorInFetch, isLoading, redirect } = this.state 
        
        /* Renders a modal while the information is being requested */
        if (isLoading){
            return (
                <div>
                    <Modal centered isOpen={true}>
                        <ModalHeader>Getting data</ModalHeader>
                        <ModalBody> 
                            <Spinner/>
                            {' '}
                            Available information is getting loaded
                        </ModalBody>
                    </Modal>
                </div>
            )
        } 
        /* Renders a modal while the information is being requested */
        if (errorInFetch){
            return(
            <Modal centered isOpen={true}>
                <ModalHeader>Error</ModalHeader>
                <ModalBody> 
                    There was an error while requesting the available information. Please try again later.
                </ModalBody>
                <ModalFooter>
                    <Link to="/" style={{ color: "#0096D6", textDecoration: "none" }}>OK!</Link>
                </ModalFooter>
            </Modal>
            )
        }
        if(redirect){
            return(
                <Navigate to="/deviceInfo"/>
            )
        }

        return(
            <div>
                <Container style={{marginTop: '2rem'}}>
                    <Row>
                        <Alert color="primary" style={{marginBottom:'1.5em'}}>
                            <div>Below you can find the available information from the devices.</div>
                            <br/>
                            <div>Select the desired one in order to further analyze it. </div>
                        </Alert>
                    </Row>
                    <Row>
                        <Button style={{marginBottom:'1.5em', backgroundColor:"#0096D6", width: "200px"}} onClick={() => window.location.reload()}>Reload information</Button>
                    </Row>
                    <Row >
                        <Col
                            className="bg-light border"
                            sm="6"
                            xs="12"
                            style={{padding: '10px', marginBottom:"25px"}}
                        >
                            <Table hover style={{color:"#0096D6"}}>
                                <thead >
                                    <tr>
                                        <th>Device's Jobs data</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {this.renderRows(this.state.jobs)}
                                </tbody>
                            </Table>
                        </Col>

                        <Col
                            className="bg-light border"
                            sm="6"
                            xs="12"
                            style={{padding: '10px', marginBottom:"25px"}}
                        >
                            <Table hover style={{color:"#0096D6"}}>
                                <thead>
                                    <tr>
                                        <th>Device's Identification data</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {this.renderRows(this.state.identification)}
                                </tbody>
                            </Table>
                        </Col>
                    </Row>
                </Container>
            </div>
        )

    }
}

export default DeviceInfoList