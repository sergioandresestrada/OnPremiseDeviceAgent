
import React from 'react';
import Heartbeat from './Heartbeat';
import Job from './Job';
import Upload from './Upload';
import { Form as FormRS, FormGroup, Input, Label} from 'reactstrap';
import '../App.css';

enum Type {
    HEARTBEAT = "HEARTBEAT",
    JOB = "JOB",
    UPLOAD = "UPLOAD"
}

const initialState = {
    type : Type.HEARTBEAT,
}


interface IJob{
    type : Type,
}

class Form extends React.Component<{}, IJob>{ 
    constructor(props : any){
        super(props);
        this.state = initialState
        this.handleChangeType = this.handleChangeType.bind(this)
    }

    
    handleChangeType(event : React.ChangeEvent<HTMLInputElement>){
        this.setState({
            type : event.target.value as Type
        })
    }

    render() {
        return(
            <div className='Form'>
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
                    {this.state.type === "HEARTBEAT" &&
                        <Heartbeat/>
                    }

                    {this.state.type === "JOB" && 
                        <Job/>
                    }

                    {this.state.type === "UPLOAD" &&
                        <Upload />
                    }
            </div>
        )
    }
}

export default Form;