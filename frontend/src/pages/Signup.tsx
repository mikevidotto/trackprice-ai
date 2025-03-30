import { useState } from "react";
import { useNavigate } from "react-router-dom";
import API from "../utils/api";
import { AuthResponse } from "../utils/types";

export default function Signup() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await API.post<AuthResponse>("/signup", { email, password });
      navigate("/login");
    } catch (err) {
      alert("Signup failed!");
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
      />
      <br></br>
      <input
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      /> 
      <br></br>
      <button type="submit">Signup</button>
    </form>
  );
}