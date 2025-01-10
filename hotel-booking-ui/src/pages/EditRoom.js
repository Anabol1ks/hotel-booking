import React, { useState, useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import Cookies from 'js-cookie'
import styles from './EditRoom.module.css'

const EditRoom = () => {
    const { id } = useParams()
    const [formData, setFormData] = useState({
			hotel_id: '',
			room_type: '',
			price: '',
			amenities: '',
			capacity: '',
		})
    const [error, setError] = useState('')
    const navigate = useNavigate()

    useEffect(() => {
        fetchRoom()
    }, [])

    const fetchRoom = async () => {
        try {
            const response = await fetch(
                process.env.REACT_APP_API_URL + `/owners/rooms?room_id=${id}`,
                {
                    headers: {
                        Authorization: `Bearer ${Cookies.get('token')}`,
                    },
                }
            )
            if (response.ok) {
                const data = await response.json()
                const room = data.find(r => r.ID === parseInt(id))
                if (room) {
                    setFormData({
                        hotel_id: room.HotelID,
                        room_type: room.RoomType,
                        price: room.Price,
                        amenities: room.Amenities,
                        capacity: room.Capacity
                    })
                }
            }
        } catch (err) {
            setError('Ошибка при загрузке данных номера')
        }
    }

    const handleSubmit = async (e) => {
        e.preventDefault()

        const submitData = {
					...formData,
					hotel_id: parseInt(formData.hotel_id),
					price: parseFloat(formData.price),
					capacity: parseInt(formData.capacity),
				}
        try {
					const response = await fetch(
						process.env.REACT_APP_API_URL + `/owners/${id}/room`,
						{
							method: 'PUT',
							headers: {
								'Content-Type': 'application/json',
								Authorization: `Bearer ${Cookies.get('token')}`,
							},
							body: JSON.stringify(submitData),
						}
					)
					if (response.ok) {
						navigate(`/owner/rooms/${id}/edit`)
					} else {
						const data = await response.json()
						setError(data.error)
					}
				} catch (err) {
					setError('Ошибка при обновлении номера')
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
            <button onClick={() => navigate('/owner/rooms')} className={styles.backButton}>
                Назад
            </button>
            <h1>Редактирование номера</h1>
            {error && <div className={styles.error}>{error}</div>}
            <form onSubmit={handleSubmit} className={styles.form}>
                <input
                    type="text"
                    name="room_type"
                    placeholder="Тип номера"
                    value={formData.room_type}
                    onChange={handleChange}
                    required
                />
                <input
                    type="number"
                    name="price"
                    placeholder="Цена за ночь"
                    value={formData.price}
                    onChange={handleChange}
                    required
                />
                <input
                    type="text"
                    name="amenities"
                    placeholder="Удобства"
                    value={formData.amenities}
                    onChange={handleChange}
                />
                <input
                    type="number"
                    name="capacity"
                    placeholder="Вместимость"
                    value={formData.capacity}
                    onChange={handleChange}
                    required
                />
                <button type="submit">Сохранить изменения</button>
            </form>
        </div>
    )
}

export default EditRoom
