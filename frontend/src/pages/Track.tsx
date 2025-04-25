import { useEffect } from "react";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import "../css/Track.css";
import API from "../utils/api";
import moment from 'moment';

export default function Track() {
    const navigate = useNavigate();
    const [url, setUrl] = useState("");
    const [competitorname, setName] = useState("");
    const [competitors, setCompetitors] = useState<any[]>([]);
    const [userData, setUserData] = useState({ email: "", firstname:"", lastname:"" })

    const submitCompetitor = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await API.post('/api/track', { "url": url, "competitor_name": competitorname })
                .then(() => {
                    getCompetitors()
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
                if (error.status == 401) {
                    console.error("token is no longer valid. logging out.")
                    API.post("/logout")
                        .then(() => {
                            localStorage.setItem("token", "")
                            navigate("/")
                        },
                            (error) => {
                                var status = error.response.status
                                console.log(status)
                            })
                } else {
                    console.error("Error fetching competitors in CompList:", error);
                }
            });
    }

    function getUserData() {
        API.get("http://localhost:8085/api/getuserdata")
            .then((response) => {
                if (response.data == null) {
                    console.error("Error: cannot retrieve user...")
                }
                if (response.data != null) {
                    setUserData(response.data.userData)
                }
            })
            .catch((error) => {
                console.error("Error retrieving user data:", error);
            });
    }

    useEffect(
        getCompetitors, // <- function that will run on every dependency update
        [] // <-- empty dependency array
    )

    useEffect(
        getUserData,
        []
    )

    return (
        <>
            <section id="dashboard">
                <aside className="dashboard-side">
                    <div id="side-data">
                        <div id="user-data">
                            <h4 id="username" >Hi, {(userData.firstname)} {(userData.lastname)}</h4>
                            <h6 id="email">{(userData.email)}</h6>
                        </div>
                        <div id="settings">
                            <p>Email Alerts:</p>
                            <label className="switch">
                                <input type="checkbox"></input>
                                <span className="slider"></span>
                            </label>
                        </div>
                    </div>
                </aside>
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
                                    {(competitors.competitor_name)} 
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
                                    {moment(competitors.last_scraped_data).format("MMM D")}
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
                            <label htmlFor="name">Name: </label>
                            <input
                                id="name"
                                type="name"
                                value={competitorname}
                                onChange={(e) => setName(e.target.value)}
                            />
                            <label htmlFor="url">Url: </label>
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
