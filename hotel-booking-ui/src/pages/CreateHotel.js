import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import Cookies from 'js-cookie'
import styles from './CreateHotel.module.css'

const CreateHotel = () => {
    const [formData, setFormData] = useState({
        name: '',
        address: '',
        description: ''
    })
    const [error, setError] = useState('')
    const navigate = useNavigate()

    const handleSubmit = async (e) => {
        e.preventDefault()
        try {
            const response = await fetch(
                process.env.REACT_APP_API_URL + '/owners/hotels',
                {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        Authorization: `Bearer ${Cookies.get('token')}`,
                    },
                    body: JSON.stringify(formData),
                }
            )
            if (response.ok) {
                navigate('/owner/hotels')
            } else {
                const data = await response.json()
                setError(data.error)
            }
        } catch (err) {
            setError('Ошибка при создании отеля')
        }
    }

    const handleChange = (e) => {
        const { name, value } = e.target
        setFormData(prev => ({
            ...prev,
            [name]: value
        }))
    }

    return (
        <div className={styles.container}>
            <button onClick={() => navigate('/owner/hotels')} className={styles.backButton}>
                Назад
            </button>
            <h1>Создание отеля</h1>
            {error && <div className={styles.error}>{error}</div>}
            <form onSubmit={handleSubmit} className={styles.form}>
                <input
                    type="text"
                    name="name"
                    placeholder="Название отеля"
                    value={formData.name}
                    onChange={handleChange}
                    required
                />
                <input
                    type="text"
                    name="address"
                    placeholder="Адрес"
                    value={formData.address}
                    onChange={handleChange}
                    required
                />
                <textarea
                    name="description"
                    placeholder="Описание"
                    value={formData.description}
                    onChange={handleChange}
                />
                <button type="submit">Создать отель</button>
            </form>
        </div>
    )
}

export default CreateHotel
