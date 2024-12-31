import React, { useState, useEffect, useCallback } from 'react'
import { useNavigate } from 'react-router-dom' // Add this import
import styles from './HotelsAndRooms.module.css'

const HotelsAndRooms = () => {
  const [hotels, setHotels] = useState([])
  const [rooms, setRooms] = useState([])
  const [error, setError] = useState('')
  const [filters, setFilters] = useState({
    minPrice: '',
    maxPrice: '',
    capacity: ''
  })

  const fetchHotels = async () => {
    try {
      const response = await fetch('http://localhost:8080/hotels')
      if (response.ok) {
        const data = await response.json()
        setHotels(data)
      }
    } catch (err) {
      setError('Ошибка при загрузке отелей')
    }
  }

  const fetchRooms = useCallback(async () => {
    try {
      const queryParams = new URLSearchParams()
      if (filters.minPrice) queryParams.append('min_price', filters.minPrice)
      if (filters.maxPrice) queryParams.append('max_price', filters.maxPrice)
      if (filters.capacity) queryParams.append('capacity', filters.capacity)

      const response = await fetch(`http://localhost:8080/rooms?${queryParams}`)
      if (response.ok) {
        const data = await response.json()
        setRooms(data)
      }
    } catch (err) {
      setError('Ошибка при загрузке номеров')
    }
  }, [filters])

  useEffect(() => {
    fetchHotels()
  }, [])

  useEffect(() => {
    fetchRooms()
  }, [fetchRooms])

  const handleFilterChange = (e) => {
    const { name, value } = e.target
    setFilters(prev => ({
      ...prev,
      [name]: value
    }))
  }

  const navigate = useNavigate()

  return (
		<div className={styles.container}>
			<div>
				<button className={styles.backButton} onClick={() => navigate('/')}>
					Вернуться на главную
				</button>
				{/* ... rest of the existing JSX ... */}
			</div>
			<h1 className={styles.title}>Отели и номера</h1>
			<div className={styles.filters}>
				<h3>Фильтры</h3>
				<div className={styles.filterGroup}>
					<input
						type='number'
						name='minPrice'
						placeholder='Мин. цена'
						value={filters.minPrice}
						onChange={handleFilterChange}
						className={styles.filterInput}
					/>
					<input
						type='number'
						name='maxPrice'
						placeholder='Макс. цена'
						value={filters.maxPrice}
						onChange={handleFilterChange}
						className={styles.filterInput}
					/>
					<input
						type='number'
						name='capacity'
						placeholder='Количество гостей'
						value={filters.capacity}
						onChange={handleFilterChange}
						className={styles.filterInput}
					/>
				</div>
			</div>

			{error && <div className={styles.error}>{error}</div>}

			<div className={styles.hotelsSection}>
				<h2>Отели</h2>
				<div className={styles.hotelsList}>
					{hotels.map(hotel => (
						<div key={hotel.ID} className={styles.hotelCard}>
							<h3>{hotel.Name}</h3>
							<p>{hotel.Address}</p>
							<p className={styles.description}>{hotel.Description}</p>
						</div>
					))}
				</div>
			</div>

			<div className={styles.roomsSection}>
				<h2>Доступные номера</h2>
				<div className={styles.roomsList}>
					{rooms.map(room => (
						<div key={room.ID} className={styles.roomCard}>
							<h3>Тип номера: {room.RoomType}</h3>
							<p>Цена: {room.Price} руб/ночь</p>
							<p>Вместимость: {room.Capacity} чел.</p>
							<p>Удобства: {room.Amenities}</p>
						</div>
					))}
				</div>
			</div>
		</div>
	)
}

export default HotelsAndRooms