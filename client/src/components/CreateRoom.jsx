import React from "react";
import { useNavigate } from "react-router-dom";

const CreateRoom = () => {
    const navigate = useNavigate();

    const create = async (e) => {
        e.preventDefault();

        const resp = await fetch("http://localhost:8080/create");
        const { room_id } = await resp.json();

        navigate(`/room/${room_id}`); // ✅ works in v6
    };

    return (
        <div>
            <button onClick={create}>Create Room</button>
        </div>
    );
};

export default CreateRoom;
