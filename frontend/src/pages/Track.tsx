import { useEffect } from "react";
import { useState } from "react";
import "../css/Track.css";
import API from "../utils/api";

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
                            <div id="competitor-box-nameurl">{competitors.competitor_url}</div>
                            <div id="competitor-box-created">{competitors.created_at}</div>
                        </div>
                    })}
                    <div className="competitor-box">
                        <form id="competitor-form" onSubmit={submitCompetitor}>
                            <input
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
