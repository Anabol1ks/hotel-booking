import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import Cookies from 'js-cookie'
import styles from './OwnerHotels.module.css'

const OwnerHotels = () => {
    const [hotels, setHotels] = useState([])
    const [error, setError] = useState('')
    const navigate = useNavigate()

    useEffect(() => {
        fetchHotels()
    }, [])

    const fetchHotels = async () => {
        try {
            const response = await fetch(
                process.env.REACT_APP_API_URL + '/owners/hotels',
                {
                    headers: {
                        Authorization: `Bearer ${Cookies.get('token')}`,
                    },
                }
            )
            if (response.ok) {
                const data = await response.json()
                setHotels(data)
            }
        } catch (err) {
            setError('Ошибка при загрузке отелей')
        }
    }

    return (
        <div className={styles.container}>
            <button onClick={() => navigate('/')} className={styles.backButton}>
                Назад
            </button>
            <h1>Мои отели</h1>
            <button 
                onClick={() => navigate('/owner/hotels/create')} 
                className={styles.createButton}
            >
                Добавить отель
            </button>
            {error && <div className={styles.error}>{error}</div>}
            <div className={styles.hotelsList}>
                {hotels.map(hotel => (
                    <div key={hotel.ID} className={styles.hotelCard}>
                        <h3>{hotel.Name}</h3>
                        <p>{hotel.Address}</p>
                        <p>{hotel.Description}</p>
                    </div>
                ))}
            </div>
        </div>
    )
}

export default OwnerHotels
