import React, { FormEvent } from "react";
import { Form as FormRS, FormGroup, Input, Label, Modal, ModalBody, ModalHeader, Spinner, Button, ModalFooter} from 'reactstrap';
import { URL } from '../utils/utils';
import Help from "./Help";

interface IHeartbeat{
    message : string,

    processingHB : boolean,
    submitOutcome : string
}

const initialState = {
    message : '',

    processingHB : false,
    submitOutcome : ''
}

class Heartbeat extends React.Component<{}, IHeartbeat>{
    constructor(props: any){
        super(props)
        this.state = initialState
        
        this.handleSubmit = this.handleSubmit.bind(this)
        this.handleChangeMessage = this.handleChangeMessage.bind(this)
    }

    handleChangeMessage(event: React.ChangeEvent<HTMLInputElement>){
        this.setState({
            message : event.target.value
        });
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
                message: this.state.message
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
        this.setState(initialState)
    }

    render(){
        return(
            <div>
                <FormRS onSubmit={this.handleSubmit}>
                    <FormGroup>
                        <Label for='heartbeatMessage'>Message to send</Label>
                        <Input onChange={this.handleChangeMessage} type="text" id="heartbeatMessage" value={this.state.message} required/>
                    </FormGroup>
                    <FormGroup>
                        <Input type="submit" value="Send Heartbeat to the printer"/>
                    </FormGroup>
                </FormRS>

                <Help 
                    message={"You can use a Heartbeat message to test whether a device is active or not.\n"+
                                "Introduce the desired message to send and click the button."} 
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
