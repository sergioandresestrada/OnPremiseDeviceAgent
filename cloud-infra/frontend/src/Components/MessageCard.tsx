import { Button, Card, CardBody, CardSubtitle, CardText, CardTitle } from "reactstrap"
import { Message } from "../utils/types"

type MessageCardProps = {
    message: Message;
  };

const MessageCard = ({ message }: MessageCardProps) => {
    
    let outlineColor = "danger"

    if (message.LastResult === "SUCCESS") {
        outlineColor = "success"
    } else if (message.LastResult === "") {
        outlineColor = "warning"
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
                <CardText>{message.AdditionalInfo}</CardText>
                <Button>Details</Button>
            </CardBody>
        </Card>
    )
}

export default MessageCard