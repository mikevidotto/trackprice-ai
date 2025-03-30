
import { useState } from "react";
import axios from 'axios'

export const client = axios.create({})

export default function Track() {
  const [url, setUrl] = useState("");
  const token = localStorage.getItem('token');
  
  const [competitors, setCompetitors] = useState<any[]>([]);

//   Array(3) [ {…}, {…}, {…} ]
// 0: Object { competitor_url: "", created_at: "2025-03-29T20:53:13.105363Z", id: 8 }
// 1: Object { competitor_url: "", created_at: "2025-03-29T20:37:52.826813Z", id: 2 }
// 2: Object { competitor_url: "", created_at: "2025-03-29T20:35:48.315199Z", id: 1 }
// length: 3
// <prototype>: Array []

  const submitCompetitor= async (e: React.FormEvent) => {
    e.preventDefault();
    console.log("URL:")
    console.log(url);
    try {
      axios.post(
                'http://localhost:8085/api/track',
                {"url": url},
                {headers: {
                    'Authorization' : 'Bearer ' + token
      }})
        .then((response) => {
                console.log("SUCCESS OR WHAT?!?")
                console.log(response.data);
                console.log(response);
                // setCompetitors(competitors => [...competitors, response.data]);
                // competitors.map(response => response.data.data[0].id)
                getCompetitors();

              },
              (error) => {
                    var status = error.response.status
                    console.log(status)
                    })

    } catch (err) {
      alert("Track failed!");
    }
  };

  const getCompetitors = async () => {
    axios.get("http://localhost:8085/api/tracked", {
      headers: {
        'Authorization' : 'Bearer ' + token
      },
    })
      .then((response) => {
        console.log(response.data);
        setCompetitors(response.data);
        console.log("COMPETITORS")
        console.log(response.data[0].id)
        console.log(response.data[0].url)
        console.log(response.data[0].created_at)
        console.log(competitors)
      })
      .catch((error) => {
        console.error("Error fetching competitors:", error);
      });
  };

  return (
    <>
    <form onSubmit={submitCompetitor}>
      <input
        type="url"
        value={url}
        onChange={(e) => setUrl(e.target.value)}
      />
      <br></br>
      <button type="submit">Track</button>
    </form>
    <div>
      {competitors.map((competitors) => {
        return <p>{competitors.id} : {competitors.competitor_url} created at: {competitors.created_at} </p>
      })}
    </div> 
    </>
  );
}
