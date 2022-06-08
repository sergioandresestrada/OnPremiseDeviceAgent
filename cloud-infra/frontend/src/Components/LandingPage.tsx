import { Link } from "react-router-dom";
import { Button, Col, Container, Row } from "reactstrap"
import '../App.css';

const LandingPage = () => {
    return(
        <div className="landing-page">
            <Container>
                <Row>
                    <Col className='landing-info custom-background' style={{marginTop: "10em", textAlign: "right"}}>
                        <h1>On-Premise Device Agent</h1>
                        <h3>Communication with devices inside a private network from a cloud environment</h3>
                    </Col>
                </Row>
                <Row>
                    <Col className='landing-info' style={{textAlign:"right"}}>
                        <div className="custom-background" style={{width: "fit-content", float: "right"}}>
                            <Button color="#0096D6" className="button-landing" style={{width: "250px", margin: "15px"}}
                                tag={Link} to="/message">
                                New Message
                            </Button>
                            <Button color="#0096D6" className="button-landing" style={{width: "250px", margin: "15px"}}
                                tag={Link} to="/devices">
                                Administrate Devices
                            </Button>
                        </div>
                    </Col>
                </Row>
            </Container>
        </div>
    )
}

export default LandingPage