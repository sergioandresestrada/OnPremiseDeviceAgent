import { AccordionHeader, AccordionItem, UncontrolledAccordion, AccordionBody, Badge } from "reactstrap";
import { Jobs } from "../utils/types"


const ShowJobs = ({Jobs}: Jobs.RootObject) => {
    return(
    <UncontrolledAccordion defaultOpen="1" open="1" style={{paddingTop:"1em"}}>
    <AccordionItem>
        <AccordionHeader targetId="1">General Information</AccordionHeader>
        <AccordionBody accordionId="1">
            <ul>
                <li key="version">{"Version: " + Jobs.Version}</li>
                <li key="date">{"Date: " + Jobs.Date}</li>
            </ul>
        </AccordionBody>
    </AccordionItem>
    <AccordionItem>
    <AccordionHeader targetId="2">
            Jobs  <Badge color="#0096D6" style={{marginInline:"1em", background:"#0096D6"}}>{Jobs.Job.length}</Badge>
        </AccordionHeader>
        <AccordionBody accordionId="2">
            <UncontrolledAccordion open="">
                {Jobs.Job.map((job, index) => {
                    return (
                        <div key={index}>
                            <AccordionItem>
                                <AccordionHeader targetId={index.toString()}>{"Job " + (index+1).toString()}</AccordionHeader>
                                <AccordionBody accordionId={index.toString()}>
                                    <ul>
                                        <li key="JobId">{"Job ID: " + job.Job_ID}</li>
                                        <li key="Status">
                                            {job.Status.Message["@hasStringResource"] === "true" ? "Status: " + job.Status.Message["#text"] : "No status provided"}
                                        </li>
                                        <li key="Links">
                                            {job.Links.Link.map((link, i) => {
                                                return(
                                                    <div key={i}>
                                                        {"Link " + (i === 0 ? "" : (index+1).toString())}
                                                        <ul>
                                                            <li key={i+"method"}>{"Method: " + link["@method"]}</li>
                                                            <li key={i+"rel"}>{"Rel: " + link["@rel"]}</li>
                                                            <li key={i+"uri"}>{"URI: " + link["@uri"]}</li>
                                                        </ul>
                                                    </div>
                                                )
                                            })}
                                        </li>
                                    </ul>
                                </AccordionBody>
                            </AccordionItem>
                        </div>
                    )
                })}
            </UncontrolledAccordion>
        </AccordionBody>
    </AccordionItem>
    <AccordionItem>
        <AccordionHeader targetId="3">
            Links  <Badge color="#0096D6" style={{marginInline:"1em", background:"#0096D6"}}>{Jobs.Links.Link.length}</Badge>
        </AccordionHeader>
        <AccordionBody accordionId="3">
            {Jobs.Links.Link.map((link, index) => {
                return (
                    <div key={index}>
                        {"Link " + (index+1).toString()}
                        <ul>
                            <li key={index+"method"}>{"Method: " + link["@method"]}</li>
                            <li key={index+"rel"}>{"Rel: " + link["@rel"]}</li>
                            <li key={index+"uri"}>{"URI: " + link["@uri"]}</li>
                        </ul>
                    </div>
                )
            }
            )}
        </AccordionBody>
    </AccordionItem>
</UncontrolledAccordion>
    )
}

export default ShowJobs