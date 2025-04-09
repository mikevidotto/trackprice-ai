import { useEffect } from "react";
import { useState } from "react";
import "../css/Track.css";
import API from "../utils/api";
import moment from 'moment';

export default function Track() {
    const [url, setUrl] = useState("");
    const [competitors, setCompetitors] = useState<any[]>([]);

    const submitCompetitor = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await API.post('/api/track', { "url": url })
                .then(() => {
                    console.log("WE ADDED A TRACKER.")
                    getCompetitors()
                    console.log("WE ARE AFTER getCompetitors")
                },
                    (error) => {
                        var status = error.response.status
                        console.log(status)
                    })
        } catch (err) {
            console.error("Failed to track competitors:", err);
        }
    };

    function getCompetitors() {
        API.get("http://localhost:8085/api/tracked")
            .then((response) => {
                if (response.data == null) {
                    console.error("Error: user has no competitors.")
                }
                if (response.data != null) {
                    setCompetitors(response.data);
                }
            })
            .catch((error) => {
                console.error("Error fetching competitors in CompList:", error);
            });
    }
    useEffect(
        getCompetitors, // <- function that will run on every dependency update
        [] // <-- empty dependency array
    )

    return (
        <>
            <section id="dashboard">
                <aside className="dashboard-side"></aside>
                <div id="dashboard-card">
                    <h4>All Competitors</h4>
                    <hr></hr>
                    <div id="filter-bar">
                        <button>Name</button>
                        <button>Created</button>
                    </div>
                    {competitors.map((competitors) => {
                        return <div className="competitor-box" key={competitors.id}>
                            <div id="box-data">
                                <div id="comp-label">
                                    Competitor
                                </div>
                                <div id="comp-url">
                                    {(competitors.competitor_url)}
                                </div>
                            </div>
                            <div id="box-data">
                                <div id="comp-label">
                                    Prices
                                </div>
                                <div id="comp-prices">
                                    <div className="comp-prices-label">
                                        Pro:
                                    </div>
                                    <div className="comp-prices">
                                        $12 
                                    </div>
                                    <div className="comp-prices-change">
                                        +8.3%
                                    </div>
                                </div>
                            </div>
                            <div id="box-data">
                                <div id="comp-label">
                                    Last Change
                                </div>
                                <div id="comp-last-change">
                                    Mar 8
                                </div>
                            </div>
                            <div id="box-data">
                                <div id="comp-label">
                                    Last Scraped
                                </div>
                                <div id="comp-last-scraped">
                                    Mar 8
                                </div>
                            </div>
                            <div id="box-data">
                                <div id="comp-label">
                                    Created
                                </div>
                                <div id="comp-created">
                                    {moment(competitors.created_at).format("MMM D")}
                                </div>
                            </div>
                        </div>
                    })}
                    <div className="submit-competitor-box">
                        <form id="competitor-form" onSubmit={submitCompetitor}>
                            <label htmlFor="url">Competitor Url: </label>
                            <input
                                id="url"
                                type="url"
                                value={url}
                                onChange={(e) => setUrl(e.target.value)}
                            />
                            <br></br>
                            <button type="submit">Track</button>
                        </form>
                    </div>
                </div>
                <aside className="dashboard-side"></aside>
            </section>
        </>
    );
}
