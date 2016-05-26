import React from "react";

var DraftHistory = React.createClass({
    render: function () {
        return (
            <div className="hist col-md-12">
                <h2>Auction History</h2>
                <div className="panel panel-default">
                    <div className="panel-body hist-container">
                        <div id="history">
                            
                        </div>
                    </div>
                </div>
            </div>
        );
    }
});


export default DraftHistory;
