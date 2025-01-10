import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import Cookies from 'js-cookie'
import styles from './OwnerRooms.module.css'

const OwnerRooms = () => {
    const [rooms, setRooms] = useState([])
    const [error, setError] = useState('')
    const navigate = useNavigate()

    useEffect(() => {
        fetchRooms()
    }, [])

    const fetchRooms = async () => {
        try {
            const response = await fetch(
                process.env.REACT_APP_API_URL + '/owners/rooms',
                {
                    headers: {
                        Authorization: `Bearer ${Cookies.get('token')}`,
                    },
                }
            )
            if (response.ok) {
                const data = await response.json()
                setRooms(data)
            }
        } catch (err) {
            setError('Ошибка при загрузке номеров')
        }
    }

    const handleDelete = async (id) => {
        if (!window.confirm('Вы уверены, что хотите удалить этот номер?')) {
            return
        }

        try {
            const response = await fetch(
                process.env.REACT_APP_API_URL + `/owners/${id}/room`,
                {
                    method: 'DELETE',
                    headers: {
                        Authorization: `Bearer ${Cookies.get('token')}`,
                    },
                }
            )
            if (response.ok) {
                fetchRooms()
            }
        } catch (err) {
            setError('Ошибка при удалении номера')
        }
    }

    return (
        <div className={styles.container}>
            <button onClick={() => navigate('/')} className={styles.backButton}>
                Назад
            </button>
            <h1>Мои номера</h1>
            <button 
                onClick={() => navigate('/owner/rooms/create')} 
                className={styles.createButton}
            >
                Добавить номер
            </button>
            {error && <div className={styles.error}>{error}</div>}
            <div className={styles.roomsList}>
                {rooms.map(room => (
                    <div key={room.ID} className={styles.roomCard}>
                        <h3>Тип номера: {room.RoomType}</h3>
                        <p>Цена: {room.Price} руб/ночь</p>
                        <p>Вместимость: {room.Capacity} чел.</p>
                        <p>Удобства: {room.Amenities}</p>
                        <div className={styles.actions}>
                            <button 
                                onClick={() => navigate(`/owner/rooms/${room.ID}/edit`)}
                                className={styles.editButton}
                            >
                                Редактировать
                            </button>
                            <button 
                                onClick={() => handleDelete(room.ID)}
                                className={styles.deleteButton}
                            >
                                Удалить
                            </button>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    )
}

export default OwnerRooms
