import React from "react";
import { Button, Card, CardBody, Collapse } from "reactstrap";

interface IHelp{
    message : string,
    opened: boolean
}

class Help extends React.Component<IHelp, IHelp>{
    constructor(props: any){
        super(props)
        this.state = props

        this.toggle = this.toggle.bind(this);
    }

    toggle() {
        this.setState({ opened: !this.state.opened });
    }

    render(){
        return(
            <div>
                <Button onClick={this.toggle} style={{ marginTop: '2rem', backgroundColor:"#0096D6"}}>Help</Button>
                <Collapse isOpen={this.state.opened}>
                <Card>
                    <CardBody>
                    {this.state.message}
                    </CardBody>
                </Card>
                </Collapse>
      </div>
        )
    }
}

export default Help;