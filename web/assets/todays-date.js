const todaysDate = () => {
    const selectTodayLink = document.querySelector('[data-module="select-todays-date"]');
    const datePicker = document.querySelector('.date-picker');

    if (selectTodayLink !== null && datePicker !== null) {
        selectTodayLink.addEventListener(
            "click",
            function (e) {
                e.preventDefault();
                datePicker.value = new Date().toJSON().slice(0, 10);
            }, false
        );
    }
};
export default todaysDate;
