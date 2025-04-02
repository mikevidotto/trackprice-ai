import Header from "../components/Header";
import { useState } from "react";
import axios from 'axios'

export const client = axios.create({})

export default function Track() {
  const [url, setUrl] = useState("");
  const token = localStorage.getItem('token');
  
  const [competitors, setCompetitors] = useState<any[]>([]);

  const submitCompetitor= async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      axios.post(
                'http://localhost:8085/api/track',
                {"url": url},
                {headers: {
                    'Authorization' : 'Bearer ' + token
      }})
        .then(() => {
                getCompetitors();
              },
              (error) => {
                    var status = error.response.status
                    console.log(status)
                    })

    getCompetitors();
    } catch (err) {
      console.error("Failed to track competitors:", err);
    }
  };

  const getCompetitors = async () => {
    axios.get("http://localhost:8085/api/tracked", {
      headers: {
        'Authorization' : 'Bearer ' + token
      },
    })
      .then((response) => {
        setCompetitors(response.data);
      })
      .catch((error) => {
        console.error("Error fetching competitors:", error);
      });
  };


  return (
    <>
    <Header />
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
        return <p key={competitors.id}>{competitors.id} : {competitors.competitor_url} created at: {competitors.created_at} </p>
                    
      })}
    </div> 
    </>
  );
}
