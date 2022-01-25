
import React, { FormEvent } from 'react';
import { Form as FormRS, FormGroup, Input, Label} from 'reactstrap';

enum Type {
    HEARTBEAT = "HEARTBEAT",
    PLACEHOLDER1 = "PLACEHOLDER1",
    PLACEHOLDER2 = "PLACEHOLDER2"
}

const initialState = {
    message : '',
    type : Type.HEARTBEAT
}

interface IJob{
    message : string,
    type : Type
}

class Form extends React.Component<{}, IJob>{ 
    constructor(props : any){
        super(props);
        this.state = initialState
        this.handleSubmit = this.handleSubmit.bind(this)
        this.handleChangeMessage = this.handleChangeMessage.bind(this)
        this.handleChangeType = this.handleChangeType.bind(this)
    }

    handleSubmit(event : FormEvent){
        event.preventDefault()
        fetch("http://localhost:12345/message",{
            method: "POST",
            headers: {
                "Content-Type": "aplication/json"
            },
            body: JSON.stringify({
                type: this.state.type,
                message: this.state.message
            })
        })
        .then(response => {
            if (response.status === 200){
                alert("New job was sent")
            } else {
                alert("There was a problem sending the new job")
            }
        })
        .catch(error => {
            alert("Se ha producido un error")
        })
        this.resetForm();
    }

    resetForm = () => {
        this.setState(initialState)
    }

    handleChangeMessage(event: React.ChangeEvent<HTMLInputElement>){
        this.setState({
            message : event.target.value
        });
    }

    handleChangeType(event : React.ChangeEvent<HTMLInputElement>){
        this.setState({
            type : event.target.value as Type
        })
    }

    render() {
        return(
            <div>
                <FormRS onSubmit={this.handleSubmit}>
                    <FormGroup>
                        <Label for='jobMessage'>Introduce the job to send</Label>
                        <Input onChange={this.handleChangeMessage} type="text" id="jobMessage" value={this.state.message} required/>
                    </FormGroup>
                    <FormGroup>
                        <Label for='jobType'>Select the type of job</Label>
                        <Input id='jobType' value={this.state.type} onChange={this.handleChangeType} type="select">
                            {Object.keys(Type).map( i => {
                                return <option key={i} value={i}>{i.charAt(0)+i.substring(1).toLowerCase()}</option>
                            })}
                        </Input>
                    </FormGroup>
                    <FormGroup>
                        <Input type="submit" value="Submit" />
                    </FormGroup>
                </FormRS>
            </div>
        )
    }
}

export default Form;