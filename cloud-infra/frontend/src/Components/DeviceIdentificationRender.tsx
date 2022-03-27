import { AccordionBody, AccordionHeader, AccordionItem, UncontrolledAccordion } from "reactstrap"
import { Identification } from "../utils/types"

const ShowIdentification = ({Identification} : Identification.RootObject) => {
    return (
        <UncontrolledAccordion defaultOpen="1" open="1" style={{paddingTop:"1em"}}>
            <AccordionItem>
                <AccordionHeader targetId="1">Regular Information</AccordionHeader>
                <AccordionBody accordionId="1">
                    <ul>
                        <li key="version">{"Version: " + Identification.Version}</li>
                        <li key="date">{"Date: " + Identification.Date}</li>
                    </ul>
                </AccordionBody>
            </AccordionItem>

            <AccordionItem>
                <AccordionHeader targetId="2">Fields</AccordionHeader>
                <AccordionBody accordionId="2">
                    {renderObject(Identification.Fields)}
                </AccordionBody>
            </AccordionItem>

            <AccordionItem>
                <AccordionHeader targetId="3">Printer Properties</AccordionHeader>
                <AccordionBody accordionId="3">
                    {renderObject(Identification.PrinterProperties)}
                </AccordionBody>
            </AccordionItem>

            <AccordionItem>
                <AccordionHeader targetId="4">Materials</AccordionHeader>
                <AccordionBody accordionId="4">
                    {renderObject(Identification.Materials)}
                </AccordionBody>
            </AccordionItem>
        </UncontrolledAccordion>
    )
}

const renderObject = (obj: Object): JSX.Element => {
    if (obj === null) return(<></>)
    return (
        <ul>
            {Object.entries(obj).map((entry, i) => {
                if (Array.isArray(entry[1])){
                    return renderArray(entry[0], entry[1])
                }
                if (typeof entry[1] === "object"){
                    return (
                        <div key={"subList"+i}>
                        <li key={i}>{entry[0]+": "}</li>
                            {renderObject(entry[1])}
                        </div>
                    )
                }
                return (<li key={i}>{entry[0]+": "+ entry[1]}</li>)
            })}
        </ul>
    )
}

const renderArray = (name: string, obj: Object): JSX.Element => {

    return(
        <UncontrolledAccordion open="" style={{paddingTop:"1em", paddingBottom:"1em"}}>
            {Object.values(obj).map((val, i) => {
                return (
                    <AccordionItem>
                        <AccordionHeader targetId={i.toString()}>{name + " " + (i+1).toString()}</AccordionHeader>
                        <AccordionBody accordionId={i.toString()}>
                            {typeof val === "object" ? renderObject(val) : val}
                        </AccordionBody>
                    </AccordionItem>
                )
            })}
        </UncontrolledAccordion>
    )
}

export default ShowIdentification