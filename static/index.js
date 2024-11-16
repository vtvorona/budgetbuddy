const openFormButton = document.getElementById("openAddPopupButton");
const addPopup = document.getElementById('addPopupForm');
const popupBg = document.getElementById('formBG');

function toggleModal() {
    const modal = document.getElementById("modal");
    
    // Проверка, если модальное окно скрыто
    if (modal.classList.contains("hidden")) {
        modal.classList.remove("hidden");
        modal.classList.add("flex");
        popupBg.classList.add("blur-md"); // Добавляем размытие при открытии
    } else {
        modal.classList.remove("flex");
        modal.classList.add("hidden");
        popupBg.classList.remove("blur"); // Убираем размытие при закрытии
    }
}

// Функция для закрытия модального окна
function closeToggleModal() {
    const modal = document.getElementById("modal");
    modal.classList.remove("flex");
    modal.classList.add("hidden");
    popupBg.classList.remove("blur-md"); // Убираем размытие
}
