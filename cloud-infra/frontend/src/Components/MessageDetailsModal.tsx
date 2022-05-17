import React from 'react'
import { Link } from 'react-router-dom'
import { Alert, Modal, ModalBody, ModalFooter, ModalHeader, Spinner } from 'reactstrap'
import { Response } from "../utils/types"
import { URL, sortResponsesByNew } from "../utils/utils"

interface PMessageDetailsModal {
    deviceUUID: string,
    messageUUID: string,
    toggleFunction: React.MouseEventHandler
}

interface IMessageDetailsModal {
    responses: Response[],

    isLoading: boolean,
    errorInFetch: boolean
}

const initialState = {
    responses: [],

    isLoading: true,
    errorInFetch: false
}

class MessageDetailsModal extends React.Component<PMessageDetailsModal, IMessageDetailsModal>{
    constructor(props: any) {
        super(props)
        this.state = initialState
    }

    componentDidMount(){
        fetch(URL + "/responses/" + this.props.deviceUUID + "/" + this.props.messageUUID)
        .then(res => res.json())
        .then(
            (result) => {
                this.setState({
                    responses: result as Response[],
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

    renderResponse(response: Response) : JSX.Element {
        if (response.Result === "SUCCESS"){
            return(
                <Alert color="success" key={response.Timestamp}>SUCCESS at {(new Date(response.Timestamp)).toLocaleString()}</Alert> 
            )
        }
        
        return (
            <Alert color="danger" key={response.Timestamp}>{response.Result} at {(new Date(response.Timestamp)).toLocaleString()}</Alert> 
        )
    }

    render() {
        const { errorInFetch, isLoading } = this.state

        if (isLoading) {
            return(
                <div>
                    <Modal centered isOpen={true}>
                        <ModalHeader>Getting data</ModalHeader>
                        <ModalBody> 
                            <Spinner/>
                            {' '}
                            Message results are getting loaded
                        </ModalBody>
                    </Modal>
                </div>
            )
        }

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

        return (
            <Modal centered isOpen={true} toggle={this.props.toggleFunction} size="xl" scrollable>
                <ModalHeader>All responses</ModalHeader>
                <ModalBody> 
                    {this.state.responses.length === 0 ?
                        <Alert color="warning">There are no responses to this message yet. Check again later.</Alert> 
                        :
                        sortResponsesByNew(this.state.responses).map((response) => {
                            return this.renderResponse(response)
                        })
                    }
                </ModalBody>
            </Modal>
        )
    }
}

export default MessageDetailsModal;