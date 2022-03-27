import { AccordionBody, AccordionHeader, AccordionItem, UncontrolledAccordion } from "reactstrap"

const GeneralInfoShow = (obj: Object) => {
    
    let a = Object.values(Object.entries(obj)[0][1])[0]
    
    return (
        <UncontrolledAccordion open="" style={{paddingTop:"1em"}}>
            {Object.entries(a as Object).map((entry, i) => {
                return(
                    <AccordionItem>
                        <AccordionHeader targetId={i.toString()}>{entry[0]}</AccordionHeader>
                       <AccordionBody accordionId={i.toString()}>{Array.isArray(entry[1]) ? renderArray(entry[0], entry[1]): renderObject(entry[1] as Object)}</AccordionBody>
                    </AccordionItem>
                )
            })}
        </UncontrolledAccordion>
    )
}

const renderObject = (obj: Object): JSX.Element => {
    if (obj === null) return(<></>)
    if (typeof obj === "number" || typeof obj === "string"){
        return (<p>{obj}</p>)
    }
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
                    <AccordionItem key={i}>
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

export default GeneralInfoShow