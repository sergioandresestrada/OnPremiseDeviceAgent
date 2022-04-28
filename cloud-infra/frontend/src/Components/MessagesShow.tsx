import React from "react";
import { Link } from "react-router-dom";
import { Alert, Col, Container, Modal, ModalBody, ModalFooter, ModalHeader, Row, Spinner } from "reactstrap";
import { Message } from "../utils/types"
import { URL, sortMessagesByNew } from "../utils/utils"
import MessageCard from "./MessageCard";
import MessageDetailsModal from "./MessageDetailsModal";

interface IMessageShow {
    messages: Message[],
    errorInFetch: boolean,
    isLoading: boolean,

    showDetailsModal: boolean,
    messageUUID: string
}

const initialState = {
    messages: [],
    errorInFetch: false,
    isLoading: true,
    
    showDetailsModal: false,
    messageUUID: ''
}

class MessagesShow extends React.Component<{}, IMessageShow> {
    constructor(props: any) {
        super(props)
        this.state = initialState

        this.showDetailsFromMessage = this.showDetailsFromMessage.bind(this)
        this.detailsModalToggle = this.detailsModalToggle.bind(this)
    }

    componentDidMount() {
        let path = window.location.pathname.split("/")
        let deviceUUID = path[path.length-1]
        
        fetch(URL + "/messages/" + deviceUUID)
        .then(res => res.json())
        .then(
            (result) => {
                this.setState({
                    messages: result as Message[],
                    isLoading: false,
                    errorInFetch: false
                })
            }
        )
        .catch(error => {
            this.setState({
                isLoading: false,
                errorInFetch: true
            })
        })
    }

    showDetailsFromMessage(messageUUID: string){
        this.setState({
            showDetailsModal: true,
            messageUUID: messageUUID
        })
    }

    detailsModalToggle(){
        this.setState({
            showDetailsModal: false,
        })
    }

    render() {
        const { errorInFetch, isLoading, showDetailsModal } = this.state

        /* Renders a modal while the messages are being requested */
        if (isLoading){
            return (
                <div>
                    <Modal centered isOpen={true}>
                        <ModalHeader>Getting data</ModalHeader>
                        <ModalBody> 
                            <Spinner/>
                            {' '}
                            Device messages are getting loaded
                        </ModalBody>
                    </Modal>
                </div>
            )
        }

        /* Renders a modal indicating that there was an error while requesting the messages */
        if (errorInFetch){
            return(
            <Modal centered isOpen={true}>
                <ModalHeader>Error</ModalHeader>
                <ModalBody> 
                    There was an error while requesting device messages. Please try again later
                </ModalBody>
                <ModalFooter>
                    <Link to="/" style={{ color: "#0096D6", textDecoration: "none" }}>OK!</Link>
                </ModalFooter>
            </Modal>
            )
        }

        let messagesByNew = sortMessagesByNew(this.state.messages)

        // Shows the modal when the details button from some message was clicked
        if (showDetailsModal) {
            let path = window.location.pathname.split("/")
            let deviceUUID = path[path.length-1]
            return(
                <div>
                <MessageDetailsModal messageUUID={this.state.messageUUID} deviceUUID={deviceUUID} toggleFunction={this.detailsModalToggle}/>
                {/* Show the messages in the background */}
                <Container style={{marginTop: "2em"}}>
                    <Row>
                        {messagesByNew.map((message) => {
                            return (
                                <Col xl="4" md="6" sm="12" key={message.MessageUUID}>
                                    <MessageCard message={message} onClickDetailsButton={() => this.showDetailsFromMessage(message.MessageUUID)} />
                                </Col>
                            )
                        })}
                    </Row>
                </Container>
            </div>
            )
            
        }

        if (this.state.messages.length === 0){
            return (
                <Container style={{marginTop: "2em"}}>
                    <Alert color="warning">No are no messages from this device! <Link to="/" style={{ color: "#0096D6", textDecoration: "none"}}>Go send some now.</Link></Alert>
                </Container>
            )
        }
        
        

        return (
            <Container style={{marginTop: "2em"}}>
                <Row>
                    {messagesByNew.map((message) => {
                        return (
                            <Col xl="4" md="6" sm="12" key={message.MessageUUID}>
                                <MessageCard message={message} onClickDetailsButton={() => this.showDetailsFromMessage(message.MessageUUID)}/>
                            </Col>
                        )
                    })}
                </Row>
            </Container>
        )
        
    }
}

export default MessagesShow;