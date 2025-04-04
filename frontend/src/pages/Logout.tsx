import API from "../utils/api";
import { useNavigate } from "react-router-dom";

export default function Logout() {
  const navigate = useNavigate();
    API.post("/logout")
        .then(() => {
            localStorage.setItem("token", "")
            navigate("/")
        }, 
        (error) => {
            var status = error.response.status
            console.log(status)
        })
    return (
    <></>
    );
}
