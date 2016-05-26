import React from "react";

var BidArea = React.createClass({
    render: function () {
        return (
            <div id="current" className="well col-md-12">
                <div className="row">
                    <div id="current_name" className="col-md-3">No one</div>
                    <div id="current_role" className="col-md-8">Role Description:You see this if something is brokend
                    </div>
                </div>
                <div className="row">
                    <div id="current_tier" className="col-md-3 text-muted">Bronze 8</div>
                    <div className="col-md-8 text-muted" href="#">Lolking Score:
                        <a className="text-info">2000</a>
                    </div>
                </div>
            </div>
        );
    }
});


export default BidArea;
