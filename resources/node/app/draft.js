import React from "react";
import $ from "jquery";
import ReactDom from "react-dom";
import BidArea from "./components/BidArea";
import DraftHistory from "./components/DraftHistory";
import CurrentPlayerPane from "./components/CurrentPlayerPane";

$(document).ready(function () {
    let di = document.createElement("div");
    ReactDom.render(
        <div className="col-lg-12">
            <BidArea />
            <DraftHistory />
            <h3 className="col-md-12">Currently Bidding On
                <span id="current_ign" className="text-info">None</span>
            </h3>

            <CurrentPlayerPane />

            <div className="col-md-6">
                <h3>Upcoming</h3>
                <div id="future" style={{overflowY: "scroll", height: "220px"}}>
                    <ul className="list-group" id="upcoming">
                    </ul>
                </div>
            </div>
            <div className="col-md-6">
                <h3>Remaining Captain Points</h3>
                <div id="auctioners" style={{overflowY: "scroll", height: "220px"}}>
                    <ul className="list-group" id="auctionerpoints">
                    </ul>
                </div>
            </div>
        </div>, di);
    document.body.appendChild(di);
});
