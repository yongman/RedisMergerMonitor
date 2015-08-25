$(function(){

var WS_HOST = 'ws://10.94.173.80:8009';
var openingObserver = Rx.Observer.create(function() { console.log('Opening socket'); });
var closingObserver = Rx.Observer.create(function() { console.log('Closing socket'); });

var stateSocket = Rx.DOM.fromWebSocket(
    WS_HOST +'/state', null, openingObserver, closingObserver);

var NewState = stateSocket.map(function(e){
    var state = JSON.parse(e.data);
    return state;
});


var Example = React.createClass({
    getInitialState: function () {
        return {
            nodes: {}
        };
    },
    componentDidMount: function() {
        var self = this;
        NewState.subscribe(
            function(obj) {
                var nodes = self.state.nodes;
                nodes = obj;
                self.setState({nodes: nodes});
            },
            function (e) {
                console.log('Error: ', e);
            },
            function () {
                console.log('Closed');
            }
        );
    },
    render: function() {
        var ths =[];
        ths.push(<th>Instance</th>);
        ths.push(<th>Uptime</th>);
        ths.push(<th>QPS</th>);
        ths.push(<th>Loading</th>);
        ths.push(<th>Role</th>);
        ths.push(<th>Master_host</th>);
        ths.push(<th>Master_port</th>);
        ths.push(<th>Master_Link</th>);
        ths.push(<th>Syncing</th>);
        ths.push(<th>My Offset</th>);
        ths.push(<th>Master Offset</th>);
        ths.push(<th>Rewriting</th>);
        var tds = _.map(this.state.nodes, function(node,key){
            var props = [];
            if (node.instantaneous_ops_per_sec > 5) {
                props.push(<td className="positive"><i className="icon checkmark"></i>{key}</td>);
            } else {
                props.push(<td><i className="icon close"></i>{key}</td>);
            }
            if (node.uptime_in_seconds > 60) {
                props.push(<td className="positive">{node.uptime_in_seconds}</td>);
            } else {
                props.push(<td className="negative">{node.uptime_in_seconds}</td>);
            }
            props.push(<td>{node.instantaneous_ops_per_sec}</td>);
            if (node.loading == '0') {
                props.push(<td className="positive">{node.loading}</td>);
            } else {
                props.push(<td className="negative">{node.loading}</td>);
            }
            props.push(<td>{node.role}</td>);
            if (node.role == 'master') {
                node.master_host = '-';
                node.master_port = '-';
                node.master_link_status = '-';
                node.master_sync_in_progress = '-';
            }
            props.push(<td>{node.master_host}</td>);
            props.push(<td>{node.master_port}</td>);
            if (node.master_link_status == 'up') {
                props.push(<td className="positive">{node.master_link_status}</td>);
            } else if (node.master_link_status == 'down'){
                props.push(<td className="negative">{node.master_link_status}</td>);
            } else {
                props.push(<td>{node.master_link_status}</td>);
            }
            if (node.master_sync_in_progress == '0') {
                props.push(<td className="positive">{node.master_sync_in_progress}</td>);
            } else if (node.master_sync_in_progress == '1') {
                props.push(<td className="negative">{node.master_sync_in_progress}</td>);
            } else {
                props.push(<td>{node.master_sync_in_progress}</td>);
            }
            props.push(<td>{node.slave_repl_offset}</td>);
            props.push(<td>{node.m_repl_offset}</td>);
            if (node.aof_rewrite_in_progress == '0') {
                props.push(<td className="positive">{node.aof_rewrite_in_progress}</td>);
            } else {
                props.push(<td className="negative">{node.aof_rewrite_in_progress}</td>);
            }
            return <tr className="center aligned">{props}</tr>;
        });

        return (
            <div className="ui vertical stripe quote segment">
            <table className="ui celled table">
                <thead>
                    <tr>{ths}</tr>
                </thead>
                <tbody>
                {tds}
                </tbody>
            </table>
            </div>
        );
    }
});

React.render(
    <Example />,
    document.getElementById('content')
);

});
