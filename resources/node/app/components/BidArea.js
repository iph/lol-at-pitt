import React from "react";

var BidArea = React.createClass({
    render: function () {
        return (
            <span>
                <div className="col-md-6">
                    <h3>Points Remaining:
                        <span id="points" className="text-primary">Connecting...</span>
                    </h3>
                </div>
                <div className="col-md-6">
                    <h3>Team:
                        <span id="team" className="text-primary">Connecting...</span>
                    </h3>
                </div>
                <form name="bet" id="bid" action="/draft/bid" method="GET" className="col-md-12">
                <div>
                    <div className="input-group">
                        <input type="text" name="amount" placeholder="Bid Amount" className="form-control"/>
                            <span className="input-group-btn">
                                <input type="submit" value="Bid" className="btn btn-success"/>
                                <input type="button" value="+1" className="btn btn-primary"/>
                                <input type="button" value="+5" className="btn btn-info"/>
                            </span>
                    </div>
                </div>
            </form>
            </span>
        );
    }
});


export default BidArea;
