import { Button, Card, CardBody, CardSubtitle, CardText, CardTitle } from "reactstrap"
import { Message } from "../utils/types"

type MessageCardProps = {
    message: Message;
    onClickDetailsButton: React.MouseEventHandler
  };

const MessageCard = ({ message, onClickDetailsButton }: MessageCardProps) => {
    
    let outlineColor = "danger"
    let lastResult = "Failure"

    if (message.LastResult === "SUCCESS") {
        outlineColor = "success"
        lastResult = "Success"

    } else if (message.LastResult === "") {
        outlineColor = "warning"
        lastResult = "Unknown"
    }

    return(
        <Card
            body
            outline
            color={outlineColor}
            style={{margin:"1em", border:"4px solid"}}
        >
            <CardBody>
                <CardTitle tag="h5">
                    {message.Type}
                </CardTitle>
                <CardSubtitle className="mb-2 text-muted" tag="h6">
                    {(new Date(message.Timestamp)).toLocaleString()}
                </CardSubtitle>
                <CardText>
                    {message.AdditionalInfo}
                    <br/>
                    Result: {lastResult}
                </CardText>
                <Button outline color={outlineColor} onClick={onClickDetailsButton}>Details</Button>
            </CardBody>
        </Card>
    )
}

export default MessageCard